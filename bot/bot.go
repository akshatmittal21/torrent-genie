package bot

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/db"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/torrent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

type msgLog struct {
	MessageID int
	Torrents  []torrent.Torrent
}

type messenger struct {
	ChatID int64
	msgLog
}

var bot *tgbotapi.BotAPI

func Init(ch chan os.Signal) error {

	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		logger.Error("Error creating bot", err)
		return err
	}

	logger.Info("Authorized on account ", bot.Self.UserName)

	// init logs
	msgLogs := make(map[int64][]msgLog)

	// init channels
	senderCh := make(chan sender, 100)
	messengerCh := make(chan messenger, 100)

	go startSender(bot, senderCh, messengerCh)

	go func(ch <-chan messenger) {
		for data := range ch {
			msgLogs[data.ChatID] = append(msgLogs[data.ChatID], msgLog{data.MessageID, data.Torrents})
		}
	}(messengerCh)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// clear cache at midnight
	c := cron.New()
	c.AddFunc("@midnight", func() {
		logger.Info("Clearing cache")
		msgLogs = make(map[int64][]msgLog)
		logger.Rotate()
	})
	c.Start()

	// start listening for updates
	messages, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.Error("Error getting updates", err)
		return err
	}

	go func() {
		defer recoverPanic(bot)

		for msg := range messages {
			if msg.Message == nil { // ignore any non-Message Updates
				continue
			}

			// Initial message
			if msg.Message.Command() == "start" {
				msg := tgbotapi.NewMessage(msg.Message.Chat.ID, constants.WELCOME_MSG)
				bot.Send(msg)
				continue
			}

			if msg.Message.Command() == "users" {
				adminID, err := strconv.ParseInt(os.Getenv("BOT_ADMIN"), 10, 64)
				if err != nil {
					logger.Error("Error getting admin id", err)
					continue
				}
				if msg.Message.Chat.ID == adminID {
					var count int64
					db.GetInstance().Find(&db.UserConfig{}).Count(&count)
					msg := tgbotapi.NewMessage(adminID, fmt.Sprintf("%d users", count))
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(msg.Message.Chat.ID, constants.INVALID_COMMAND)
					bot.Send(msg)
				}
				continue

			}
			if msg.Message.IsCommand() {
				go sendCommandResponse(msg, msg.Message.Command(), senderCh)
				continue
			}

			isReply, err := strconv.Atoi(msg.Message.Text)
			if err == nil {
				go sendMagnet(msg, msgLogs, isReply, senderCh)
			} else {
				go sendTorrents(msg, senderCh)
			}

			// Update user in DB
			chat := msg.Message.Chat
			user := db.UserConfig{UserID: chat.ID, FirstName: chat.FirstName, LastName: chat.LastName, UserName: chat.UserName}
			if db.GetInstance().Model(user).Where("user_id = ?", chat.ID).Updates(&user).RowsAffected == 0 {
				db.GetInstance().Create(&user)
			}
		}
	}()

	signalType := <-ch
	fmt.Println("Exit command received, Exiting...")
	fmt.Println("Received signal type : ", signalType)

	shutdown(bot)
	c.Stop()
	return nil
}

// panic recover
func recoverPanic(bot *tgbotapi.BotAPI) {
	if err := recover(); err != nil {
		logger.Error("panic occurred:", err)
		if bot != nil {
			notifyAdmin(bot, "!!!Panic Occured!!!")
		}
	}
}

func notifyAdmin(bot *tgbotapi.BotAPI, message string) {
	adminID, err := strconv.ParseInt(os.Getenv("BOT_ADMIN"), 10, 64)
	if err != nil {
		logger.Error("Error getting admin id", err)
		return
	}
	msg := tgbotapi.NewMessage(adminID, message)
	bot.Send(msg)
}

func shutdown(bot *tgbotapi.BotAPI) {
	// Gracefully shutting down server
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// sending notification to admin
	adminCh := make(chan struct{})
	go func() {
		notifyAdmin(bot, "!!! Server shutdown !!!")
		adminCh <- struct{}{}
	}()

	select {
	case <-adminCh:
		logger.Info("Admin notified")

	case <-tc.Done():
		logger.Info("Timeout reached")
	}
	bot.StopReceivingUpdates()
}
