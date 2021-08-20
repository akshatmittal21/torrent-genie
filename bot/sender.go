package bot

import (
	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/magnet"
	"github.com/akshatmittal21/torrent-genie/torrent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type sender struct {
	ChatID    int64
	Type      constants.MessageType
	MsgConfig tgbotapi.MessageConfig
	Torrents  []torrent.Torrent
}

// Default sender
func startSender(bot *tgbotapi.BotAPI, senderCh <-chan sender, messengerCh chan<- messenger) {
	defer recoverPanic(bot)

	for data := range senderCh {
		var err error
		var update tgbotapi.Message
		switch data.Type {
		case constants.Magnet:
			_, err = bot.Send(data.MsgConfig)
		case constants.Torrent:
			update, err = bot.Send(data.MsgConfig)
			if err == nil {
				messengerCh <- messenger{ChatID: data.ChatID, msgLog: msgLog{MessageID: update.MessageID, Torrents: data.Torrents}}
			}

		}
		if err != nil {
			logger.Error(err)
			replyMsg := tgbotapi.NewMessage(data.ChatID, constants.SOMETHING_WENT_WRONG)
			bot.Send(replyMsg)
		}
	}
}

// Sending magnet link
func sendMagnet(msg tgbotapi.Update, msgLogs map[int64][]msgLog, replyNo int, senderCh chan<- sender) {
	defer recoverPanic(bot)
	var msgID int
	var torrents []torrent.Torrent
	var replyMsg tgbotapi.MessageConfig

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
	if replyNo > 0 && len(torrents) > 0 && len(torrents) >= replyNo {
		tor := torrents[replyNo-1+20]
		magnetLink := magnet.GetLink(tor.InfoHash, tor.Name)

		msgstring := tor.Name + "\n\n" + "Copy the magnet below" + "\n\n" + "`" + magnetLink + "`"
		replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, msgstring)
		replyMsg.ParseMode = tgbotapi.ModeMarkdown
		replyMsg.ReplyToMessageID = msg.Message.MessageID

	} else {
		replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, constants.INVALID_REPLY)

	}
	senderCh <- sender{ChatID: msg.Message.Chat.ID, Type: constants.Magnet, MsgConfig: replyMsg}
}

// Sending torrents
func sendTorrents(msg tgbotapi.Update, senderCh chan<- sender) {
	defer recoverPanic(bot)
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
	replyMsg := tgbotapi.NewMessage(msg.Message.Chat.ID, torrentResp)
	replyMsg.ReplyToMessageID = msg.Message.MessageID
	senderCh <- sender{ChatID: msg.Message.Chat.ID, Type: constants.Torrent, MsgConfig: replyMsg, Torrents: torrents}

}
