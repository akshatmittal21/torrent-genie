package telegram

import (
	"fmt"
	"strings"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Default sender
func (b *Bot) startSender(senderCh <-chan sender, messengerCh chan<- messenger) {
	defer b.recoverPanic()

	for data := range senderCh {
		var err error
		var update tgbotapi.Message
		switch data.Type {
		case constants.Magnet:
			_, err = b.tgbot.Send(data.MsgConfig)
		case constants.Torrent:
			update, err = b.tgbot.Send(data.MsgConfig)
			if err == nil {
				messengerCh <- messenger{ChatID: data.ChatID, msgLog: msgLog{MessageID: update.MessageID, Torrents: data.Torrents}}
			}
		}
		if err != nil {
			b.log.Error(err)
			replyMsg := tgbotapi.NewMessage(data.ChatID, constants.SOMETHING_WENT_WRONG)
			_, err = b.tgbot.Send(replyMsg)
			if err != nil {
				b.log.Error(fmt.Errorf("Error sending message: %s", err))
			}
		}
	}
}

// Sending magnet link
func (b *Bot) sendMagnet(msg tgbotapi.Update, replyNo int, senderCh chan<- sender) {
	defer b.recoverPanic()
	var (
		msgID    int
		torrents []dto.Torrent
		replyMsg tgbotapi.MessageConfig
	)
	if msg.Message.ReplyToMessage != nil {
		msgID = msg.Message.ReplyToMessage.MessageID
	}
	msglog := b.msgLog[msg.Message.Chat.ID]
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
		tor := torrents[replyNo-1]

		torrentData := b.magnet.GetFile(tor.InfoHash)

		if torrentData != nil {
			file := tgbotapi.FileBytes{Name: tor.Name + ".torrent", Bytes: torrentData}
			newMsg := tgbotapi.NewDocumentUpload(msg.Message.Chat.ID, file)

			magnetLink := b.magnet.GetLink(tor.InfoHash, tor.Name)
			msgstring := "Copy the magnet below or download the torrent file" + "\n\n" + "`" + magnetLink + "`"

			newMsg.ReplyToMessageID = msg.Message.MessageID
			newMsg.Caption = msgstring
			newMsg.ParseMode = tgbotapi.ModeMarkdown
			senderCh <- sender{ChatID: msg.Message.Chat.ID, Type: constants.Magnet, MsgConfig: newMsg}
			return
		} else {
			magnetLink := b.magnet.GetLink(tor.InfoHash, tor.Name)
			msgstring := tor.Name + "\n\n" + "Copy the magnet below" + "\n\n" + "`" + magnetLink + "`"
			replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, msgstring)
			replyMsg.ParseMode = tgbotapi.ModeMarkdown
			replyMsg.ReplyToMessageID = msg.Message.MessageID
		}

	} else {
		replyMsg = tgbotapi.NewMessage(msg.Message.Chat.ID, constants.INVALID_REPLY)

	}
	senderCh <- sender{ChatID: msg.Message.Chat.ID, Type: constants.Magnet, MsgConfig: replyMsg}
}

// Sending torrents
func (b *Bot) sendTorrents(msg tgbotapi.Update, senderCh chan<- sender) {
	defer b.recoverPanic()
	var torrentResp string
	torrents, err := b.torrent.GetTorrents(msg.Message.Text)
	if err != nil {
		torrentResp = constants.SOMETHING_WENT_WRONG
	} else if len(torrents) > 0 {
		if torrents[0].ID == "0" {
			torrentResp = constants.NO_RESULTS
		} else {
			torrentResp = dto.CreateResponse(torrents)
			torrentResp = torrentResp + "\nReply with option number to get torrent.."
		}
	} else {
		torrentResp = constants.NO_RESULTS
	}
	replyMsg := tgbotapi.NewMessage(msg.Message.Chat.ID, torrentResp)
	replyMsg.ReplyToMessageID = msg.Message.MessageID
	senderCh <- sender{ChatID: msg.Message.Chat.ID, Type: constants.Torrent, MsgConfig: replyMsg, Torrents: torrents}

}

// Sending torrents
func (b *Bot) sendCommandResponse(msg tgbotapi.Update, command string, senderCh chan<- sender) {
	defer b.recoverPanic()
	var torrentResp string
	var torrents []dto.Torrent
	torrents, err := b.torrent.GetRecentTorrents(constants.Command(command))
	if err != nil {
		torrentResp = constants.SOMETHING_WENT_WRONG
	} else if len(torrents) > 0 {
		if torrents[0].ID == "0" {
			torrentResp = constants.NO_RESULTS
		} else {
			torrentResp = dto.CreateResponse(torrents)
			torrentResp = torrentResp + "\nReply with option number to get torrent.."
		}
	} else {
		torrentResp = constants.NO_RESULTS
	}
	replyMsg := tgbotapi.NewMessage(msg.Message.Chat.ID, torrentResp)
	replyMsg.ReplyToMessageID = msg.Message.MessageID
	senderCh <- sender{ChatID: msg.Message.Chat.ID, Type: constants.Torrent, MsgConfig: replyMsg, Torrents: torrents}

}

// sendRecommendMsg
func (b *Bot) sendRecommendMsg() {
	if constants.IsRecommendationOn {
		defer b.recoverPanic()
		b.log.Info("Sending Recommendation Message")
		users := b.db.GetAllUsers()
		for _, user := range users {
			msg := strings.Replace(constants.RecommendMsg, "${NAME}", user.FirstName, 1)
			msgConf := tgbotapi.NewMessage(user.UserID, msg)
			_, err := b.tgbot.Send(msgConf)
			if err != nil {
				b.log.Error(fmt.Errorf("error sending message: %s", err))
			}
		}
	}
}
