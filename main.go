package main

import (
	"log"
	"os"
	"strconv"

	"github.com/akshatmittal21/torrent-genie/magnet"
	"github.com/akshatmittal21/torrent-genie/torrent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MessageLog struct {
	MessageID int
	Torrents  []torrent.Torrent
}

func main() {
	// fmt.Println(torrent.GetTorrents("harry potter"))
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	msgLogs := make(map[int64][]MessageLog)
	messages, err := bot.GetUpdatesChan(u)
	for msg := range messages {
		if msg.Message == nil { // ignore any non-Message Updates
			continue
		}

		var replyMsg tgbotapi.MessageConfig
		if msg.Message.Command() == "start" {
			newMsg := "Welcome to TorrentGenie \n\n Search torrents by sending a name"
			msg := tgbotapi.NewMessage(msg.Message.Chat.ID, newMsg)
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

				msgstring := "Torrent: " + tor.Name + "\n\n" + "`" + magnetLink + "`"
				replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, msgstring)
				replyMsg.ParseMode = tgbotapi.ModeMarkdown
				replyMsg.ReplyToMessageID = msg.Message.MessageID

			} else {
				replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, "Invalid Reply - Try again")

			}

			_, err := bot.Send(replyMsg)
			if err != nil {
				log.Println(err)
			}

		} else {
			var torrentResp string
			torrents := torrent.GetTorrents(msg.Message.Text)
			if len(torrents) > 0 {
				if torrents[0].ID == "0" {
					torrentResp = "No torrents found"
				} else {
					torrentResp = torrent.CreateResponse(torrents)
					torrentResp = torrentResp + "\nReply with option number to get magnet link.."
				}
			} else {
				torrentResp = "No torrents found"
			}
			replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, torrentResp)
			replyMsg.ReplyToMessageID = msg.Message.MessageID
			update, err := bot.Send(replyMsg)
			if err != nil {
				log.Println(err)
			}
			msgLogs[msg.Message.Chat.ID] = append(msgLogs[msg.Message.Chat.ID], MessageLog{update.MessageID, torrents})
		}
	}
}
