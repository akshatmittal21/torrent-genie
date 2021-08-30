package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/akshatmittal21/torrent-genie/bot"
)

func main() {

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGTERM)

	// init bot
	bot.Init(sigCh)

}
