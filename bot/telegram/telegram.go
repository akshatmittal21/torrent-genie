package telegram

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/dto"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/magnet"
	"github.com/akshatmittal21/torrent-genie/store"
	"github.com/akshatmittal21/torrent-genie/torrent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

func NewBot(token, adminID string, log logger.Logger, db store.Database) (*Bot, error) {
	// Initiate bot
	tgBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		err = fmt.Errorf("Error creating bot %w", err)
		log.Error(err)
		return nil, err
	}
	return &Bot{
		adminID:       adminID,
		log:           log,
		db:            db,
		torrent:       torrent.NewServer(log),
		magnet:        magnet.NewServer(log),
		tgbot:         tgBot,
		isServerAlive: true,
		msgLog:        make(map[int64][]msgLog),
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {

	var err error
	// init channels
	senderCh := make(chan sender, 100)
	messengerCh := make(chan messenger, 100)
	b.log.Info("Authorized on account ", b.tgbot.Self.UserName)

	go b.startSender(senderCh, messengerCh)
	go b.checkTorrentServerStatus()
	go func(ch <-chan messenger) {
		for data := range ch {
			b.msgLog[data.ChatID] = append(b.msgLog[data.ChatID], msgLog{data.MessageID, data.Torrents})
		}
	}(messengerCh)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	scheduler := b.invokeScheduler()
	defer scheduler.Stop()

	// start listening for updates
	messages, err := b.tgbot.GetUpdatesChan(u)
	if err != nil {
		b.log.Error("Error getting updates", err)
		return err
	}

	go func() {
		defer b.recoverPanic()

		for msg := range messages {
			if msg.Message == nil { // ignore any non-Message Updates
				continue
			}

			if msg.Message.Command() == "start" {
				msg := tgbotapi.NewMessage(msg.Message.Chat.ID, constants.WELCOME_MSG)
				_, err = b.tgbot.Send(msg)
				if err != nil {
					b.log.Error("Error sending welcome message", err)
				}
				continue
			}

			if msg.Message.Command() == "users" {
				adminID, err := strconv.ParseInt(os.Getenv("BOT_ADMIN"), 10, 64)
				if err != nil {
					b.log.Error("Error getting admin id", err)
					continue
				}
				if msg.Message.Chat.ID == adminID {
					count := b.db.GetUserCount()
					msg := tgbotapi.NewMessage(adminID, fmt.Sprintf("%d users", count))
					_, err = b.tgbot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(msg.Message.Chat.ID, constants.INVALID_COMMAND)
					_, err = b.tgbot.Send(msg)
				}
				if err != nil {
					b.log.Error("Error sending message %w", err)
				}
				continue

			}

			// Checking server status
			if !b.isServerAlive {
				msg := tgbotapi.NewMessage(msg.Message.Chat.ID, constants.SERVER_DOWN)
				_, err := b.tgbot.Send(msg)
				if err != nil {
					b.log.Error("Error sending message %w", err)
				}
				continue
			}

			if msg.Message.Command() == "togglerecommend" {
				adminID, err := strconv.ParseInt(os.Getenv("BOT_ADMIN"), 10, 64)
				if err != nil {
					b.log.Error("Error getting admin id", err)
					continue
				}
				if msg.Message.Chat.ID == adminID {
					constants.IsRecommendationOn = !constants.IsRecommendationOn
					msg := tgbotapi.NewMessage(adminID, fmt.Sprintf("Recommendation is %s", strconv.FormatBool(constants.IsRecommendationOn)))
					_, err := b.tgbot.Send(msg)
					if err != nil {
						b.log.Error("Error sending message %w", err)
					}
				}
				continue
			}
			if msg.Message.IsCommand() {
				go b.sendCommandResponse(msg, msg.Message.Command(), senderCh)
				continue
			}

			isReply, err := strconv.Atoi(msg.Message.Text)
			if err == nil {
				go b.sendMagnet(msg, isReply, senderCh)
			} else {
				go b.sendTorrents(msg, senderCh)
			}

			// Update user in DB
			chat := msg.Message.Chat
			user := dto.DBUser{FirstName: chat.FirstName, LastName: chat.LastName, UserName: chat.UserName}
			b.db.UpsertUser(chat.ID, user)
		}
	}()
	<-ctx.Done()
	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	// sending notification to admin
	adminCh := make(chan struct{})
	go func() {
		err := b.notifyAdmin("!!! Server shutdown !!!")
		if err != nil {
			b.log.Error("Error sending message to admin", err)
		}
		adminCh <- struct{}{}
	}()

	select {
	case <-adminCh:
		b.log.Info("Admin notified")

	case <-ctx.Done():
		b.log.Info("Timeout reached")
	}
	b.tgbot.StopReceivingUpdates()
	return nil
}

// panic recover
func (b *Bot) recoverPanic() {
	if err := recover(); err != nil {
		b.log.Error("panic occurred:", err)
		err = b.notifyAdmin("!!! Panic occured !!!")
		if err != nil {
			b.log.Error("Error sending message to admin: %w", err)
		}
	}
}

func (b *Bot) notifyAdmin(message string) error {
	adminID, err := strconv.ParseInt(b.adminID, 10, 64)
	if err != nil {
		err = fmt.Errorf("Error getting admin id: %w", err)
		return err
	}

	msg := tgbotapi.NewMessage(adminID, message)
	_, err = b.tgbot.Send(msg)
	if err != nil {
		err = fmt.Errorf("Error notifying admin: %w", err)
		return err
	}
	return nil
}

func (b *Bot) checkTorrentServerStatus() {
	for range time.Tick(time.Second * constants.PingTimeout) {
		b.isServerAlive = b.torrent.IsServerAlive()
	}
}

func (b *Bot) resetInMemoryData() {
	b.log.Info("Clearing cache")
	b.msgLog = make(map[int64][]msgLog)
	b.log.Info("Rotating Logs")
	err := b.log.Rotate()
	b.log.Error("Error rotating logs %w", err)
}

func (b *Bot) invokeScheduler() *cron.Cron {
	// clear cache at midnight
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		b.log.Error("cron error", err)
		err = b.notifyAdmin("!!! cron error !!!")
		if err != nil {
			b.log.Error("Error notifying admin: %w", err)
		}
	}
	c := cron.NewWithLocation(loc)
	err = c.AddFunc("@midnight", b.resetInMemoryData)
	if err != nil {
		b.log.Error("cron error", err)
	}

	// invoke recommender
	err = c.AddFunc("0 0 17 * * *", b.sendRecommendMsg)
	if err != nil {
		b.log.Error("cron error", err)
	}
	c.Start()
	return c
}
