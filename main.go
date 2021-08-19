package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/akshatmittal21/torrent-genie/bot"
	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/logger"
	"go.uber.org/zap/zapcore"
)

func main() {

	logger.InitLogger(constants.LogPath, zapcore.InfoLevel)
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGTERM)

	// init bot
	bot.Init(sigCh)
}
