package bot

import (
	"context"
	"fmt"
	"os"

	"github.com/akshatmittal21/torrent-genie/bot/telegram"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/store"
)

type Bot interface {
	Start(context.Context) error
	Stop(context.Context) error
}

func New(log logger.Logger, db store.Database) (Bot, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN env variable not set")
	}
	adminID := os.Getenv("BOT_ADMIN")
	if adminID == "" {
		return nil, fmt.Errorf("BOT_ADMIN env variable not set")
	}
	return telegram.NewBot(botToken, adminID, log, db)
}
