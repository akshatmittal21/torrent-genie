package server

import (
	"context"
	"fmt"
	oslog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akshatmittal21/torrent-genie/bot"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/store"
)

type Server struct {
}

func (s *Server) Start() {

	log := initLog("./logs/system/log", "development")
	db := initDB("./db/users.db", log)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	tgbot, err := bot.New(log, db)
	if err != nil {
		oslog.Fatal(err)
		os.Exit(1)
	}
	ctx := context.Background()
	go func() {
		if err := tgbot.Start(ctx); err != nil {
			oslog.Fatal(err)
			os.Exit(1)
		}
	}()
	<-sigCh
	tc, stop := context.WithTimeout(ctx, 10*time.Second)
	defer stop()
	err = tgbot.Stop(tc)
	if err != nil {
		oslog.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("Shutting down...")
}

func initLog(logPath, env string) logger.Logger {
	logLevel := logger.InfoLevel
	log, err := logger.Init(logPath, logLevel)
	if err != nil {
		oslog.Fatal(err)
	}
	return log
}

func initDB(filePath string, log logger.Logger) store.Database {
	db, err := store.New(filePath, log)
	if err != nil {
		oslog.Fatal(err)
	}
	return db
}
