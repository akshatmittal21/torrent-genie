package telegram

import (
	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/dto"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/magnet"
	"github.com/akshatmittal21/torrent-genie/store"
	"github.com/akshatmittal21/torrent-genie/torrent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	adminID       string
	log           logger.Logger
	torrent       torrent.Server
	magnet        magnet.Magnet
	db            store.Database
	tgbot         *tgbotapi.BotAPI
	msgLog        map[int64][]msgLog
	isServerAlive bool
}

type msgLog struct {
	MessageID int
	Torrents  []dto.Torrent
}

type messenger struct {
	ChatID int64
	msgLog
}

type sender struct {
	ChatID    int64
	Type      constants.MessageType
	MsgConfig tgbotapi.Chattable
	Torrents  []dto.Torrent
}
