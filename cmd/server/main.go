package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/myrunes/backend/internal/caching"
	"github.com/myrunes/backend/internal/config"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/ddragon"
	"github.com/myrunes/backend/internal/logger"
	"github.com/myrunes/backend/internal/mailserver"
	"github.com/myrunes/backend/internal/webserver"
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

	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.MongoDB.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		cfg.MongoDB.Port = v
	}
	if v := os.Getenv("DB_USERNAME"); v != "" {
		cfg.MongoDB.Username = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.MongoDB.Password = v
	}
	if v := os.Getenv("DB_AUTHDB"); v != "" {
		cfg.MongoDB.AuthDB = v
	}
	if v := os.Getenv("DB_DATADB"); v != "" {
		cfg.MongoDB.DataDB = v
	}
	if v := strings.ToLower(os.Getenv("TLS_ENABLE")); v == "true" || v == "t" || v == "1" {
		cfg.WebServer.TLS.Enabled = true
	}
	if v := os.Getenv("TLS_KEY"); v != "" {
		cfg.WebServer.TLS.Key = v
	}
	if v := os.Getenv("TLS_CERT"); v != "" {
		cfg.WebServer.TLS.Cert = v
	}

	logger.Info("DDRAGON :: initialization")
	if ddragon.DDragonInstance, err = ddragon.Fetch("latest"); err != nil {
		logger.Fatal("DDRAGON :: failed polling data from ddragon: %s", err.Error())
	}
	logger.Info("DDRAGON :: initialized")

	db := new(database.MongoDB)
	logger.Info("DATABASE :: initialization")
	if err = db.Connect(cfg.MongoDB); err != nil {
		logger.Fatal("DATABASE :: failed establishing connection to database: %s", err.Error())
	}
	defer func() {
		logger.Info("DATABASE :: teardown")
		db.Close()
	}()

	logger.Info("MAILSERVER :: initialization")
	ms, err := mailserver.NewMailServer(cfg.MailServer, "noreply@myrunes.com", "myrunes")
	if err != nil {
		logger.Fatal("MAILSERVER :: failed connecting to mail account: %s", err.Error())
	}
	logger.Info("MAILSERVER :: started")

	var cache caching.CacheMiddleware
	if cfg.Redis != nil && cfg.Redis.Enabled {
		cache = caching.NewRedis(cfg.Redis)
	} else {
		cache = caching.NewInternal()
	}
	cache.SetDatabase(db)

	logger.Info("WEBSERVER :: initialization")
	ws, err := webserver.NewWebServer(db, cache, ms, cfg.WebServer)
	if err != nil {
		logger.Fatal("WEBSERVER :: failed creating web server: %s", err.Error())
	}
	go func() {
		if err := ws.ListenAndServeBlocking(); err != nil {
			logger.Fatal("WEBSERVER :: failed starting web server: %s", err.Error())
		}
	}()
	logger.Info("WEBSERVER :: started")

	// Lifecycle Timer was used to clean up expierd
	// sessions which is no more necessary after
	// implementation of JWT tokens.
	// Just keeping this here in case of this may
	// be needed some time later.
	// lct := lifecycletimer.New(5 * time.Minute).
	// 	Handle(func() {
	// 	}).
	// 	Start()
	// defer lct.Stop()
	// logger.Info("LIFECYCLETIMER :: started")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
