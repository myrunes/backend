package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/zekroTJA/lol-runes/internal/config"
	"github.com/zekroTJA/lol-runes/internal/database"
	"github.com/zekroTJA/lol-runes/internal/logger"
	"github.com/zekroTJA/lol-runes/internal/webserver"
)

var (
	flagConfig = flag.String("c", "config.yml", "config file location")
)

func main() {
	flag.Parse()

	logger.Setup(`%{color}▶  %{level:.4s} %{id:03d}%{color:reset} %{message}`, 5)

	logger.Info("CONFIG :: initialization")
	cfg, err := config.Open(*flagConfig)
	if err != nil {
		logger.Fatal("CONFIG :: failed creating or opening config: %s", err.Error())
	}
	if cfg == nil {
		logger.Info("CONFIG :: config file was created at '%s'. Set your config values and restart.", *flagConfig)
		return
	}

	db := new(database.MongoDB)
	logger.Info("DATABASE :: initialization")
	if err = db.Connect(cfg.MongoDB); err != nil {
		logger.Fatal("DATABASE :: failed establishing connection to database: %s", err.Error())
	}
	defer func() {
		logger.Info("DATABASE :: teardown")
		db.Close()
	}()

	logger.Info("WEBSERVER :: initialization")
	ws := webserver.NewWebServer(db, cfg.WebServer)
	go func() {
		if err := ws.ListenAndServeBlocking(); err != nil {
			logger.Fatal("WEBSERVER :: failed starting web server: %s", err.Error())
		}
	}()
	logger.Info("WEBSERVER :: started")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
