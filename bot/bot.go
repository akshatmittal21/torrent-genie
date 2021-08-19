package bot

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/magnet"
	"github.com/akshatmittal21/torrent-genie/torrent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

type MessageLog struct {
	MessageID int
	Torrents  []torrent.Torrent
}

func Init(ch chan os.Signal) error {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		logger.Error("Error creating bot", err)
		return err
	}

	logger.Info("Authorized on account ", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	msgLogs := make(map[int64][]MessageLog)

	// clear cache at midnight
	c := cron.New()
	c.AddFunc("@midnight", func() {
		logger.Info("Clearing cache")
		msgLogs = make(map[int64][]MessageLog)
		logger.Rotate()
	})
	c.Start()
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
			var replyMsg tgbotapi.MessageConfig
			if msg.Message.Command() == "start" {
				msg := tgbotapi.NewMessage(msg.Message.Chat.ID, constants.WELCOME_MSG)
				bot.Send(msg)
				continue

			}
			isReply, err := strconv.Atoi(msg.Message.Text)
			if err == nil {
				var msgID int
				var torrents []torrent.Torrent
				if msg.Message.ReplyToMessage != nil {
					msgID = msg.Message.ReplyToMessage.MessageID
				}
				msglog := msgLogs[msg.Message.Chat.ID]
				if len(msglog) > 0 {
					if msgID == 0 {
						torrents = msglog[len(msglog)-1].Torrents
					} else {
						for _, m := range msglog {
							if m.MessageID == msgID {
								torrents = m.Torrents
								break
							}
						}
					}
				}
				if isReply > 0 && len(torrents) > 0 && len(torrents) >= isReply {
					tor := torrents[isReply-1]
					magnetLink := magnet.GetLink(tor.InfoHash, tor.Name)

					msgstring := "Copy the magnet below for: " + tor.Name + "\n\n" + "`" + magnetLink + "`"
					replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, msgstring)
					replyMsg.ParseMode = tgbotapi.ModeMarkdown
					replyMsg.ReplyToMessageID = msg.Message.MessageID

				} else {
					replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, constants.INVALID_REPLY)

				}

				_, err := bot.Send(replyMsg)
				if err != nil {
					logger.Error(err)
					replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, constants.SOMETHING_WENT_WRONG)
					bot.Send(replyMsg)
				}

			} else {
				var torrentResp string
				torrents := torrent.GetTorrents(msg.Message.Text)
				if len(torrents) > 0 {
					if torrents[0].ID == "0" {
						torrentResp = constants.NO_RESULTS
					} else {
						torrentResp = torrent.CreateResponse(torrents)
						torrentResp = torrentResp + "\nReply with option number to get magnet link.."
					}
				} else {
					torrentResp = constants.NO_RESULTS
				}
				replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, torrentResp)
				replyMsg.ReplyToMessageID = msg.Message.MessageID
				update, err := bot.Send(replyMsg)
				if err != nil {
					logger.Error(err)
					replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, constants.SOMETHING_WENT_WRONG)
					bot.Send(replyMsg)
				}
				msgLogs[msg.Message.Chat.ID] = append(msgLogs[msg.Message.Chat.ID], MessageLog{update.MessageID, torrents})
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
	}
	notifyAdmin(bot, "!!!Panic Occured!!!")
}

func notifyAdmin(bot *tgbotapi.BotAPI, message string) {
	adminID, err := strconv.Atoi(os.Getenv("BOT_ADMIN"))
	if err != nil {
		logger.Error("Error getting admin id", err)
		return
	}
	msg := tgbotapi.NewMessage(int64(adminID), message)
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
