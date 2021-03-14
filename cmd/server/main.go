package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/myrunes/backend/pkg/ddragon"
	"github.com/myrunes/backend/pkg/lifecycletimer"

	"github.com/myrunes/backend/internal/assets"
	"github.com/myrunes/backend/internal/caching"
	"github.com/myrunes/backend/internal/config"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/logger"
	"github.com/myrunes/backend/internal/mailserver"
	"github.com/myrunes/backend/internal/storage"
	"github.com/myrunes/backend/internal/webserver"
)

var (
	flagConfig    = flag.String("c", "config.yml", "config file location")
	flagSkipFetch = flag.Bool("skipFetch", false, "skip avatar asset fetching")
)

func initStorage(c *config.Main) (st storage.Middleware, err error) {
	var cfg interface{}

	switch c.Storage.Typ {
	case "file", "fs":
		st = new(storage.File)
		cfg = *c.Storage.File
	case "minio", "s3":
		st = new(storage.Minio)
		cfg = *c.Storage.Minio
	default:
		return nil, errors.New("invalid storage type")
	}

	if cfg == nil {
		return nil, errors.New("invalid storage config")
	}

	err = st.Init(cfg)

	return
}

func fetchAssets(a *assets.AvatarHandler) error {
	if *flagSkipFetch {
		return nil
	}

	cChamps := make(chan string)
	cError := make(chan error)

	go a.FetchAll(cChamps, cError)

	go func() {
		for _, c := range ddragon.DDragonInstance.Champions {
			cChamps <- c.UID
		}
		close(cChamps)
	}()

	for err := range cError {
		if err != nil {
			return err
		}
	}

	return nil
}

func refetch(a *assets.AvatarHandler) {
	var err error

	logger.Info("DDRAGON :: refetch")
	if ddragon.DDragonInstance, err = ddragon.Fetch("latest"); err != nil {
		logger.Error("DDRAGON :: failed polling data from ddragon: %s", err.Error())
	}

	logger.Info("ASSETHANDLER :: refetch")
	if err = fetchAssets(a); err != nil {
		logger.Fatal("ASSETHANDLER :: failed fetching assets: %s", err.Error())
	}
}

func cleanupExpiredRefreshTokens(db database.Middleware) {
	n, err := db.CleanupExpiredTokens()
	if err != nil {
		logger.Error("DATABASE :: failed cleaning up expired refresh tokens: %s", err.Error())
	} else {
		logger.Info("AUTH :: cleaned %d expired refresh tokens", n)
	}
}

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

	logger.Info("STORAGE :: initialization")
	st, err := initStorage(cfg)
	if err != nil {
		logger.Fatal("STORAGE :: failed initializing storage: %s", err.Error())
	}

	logger.Info("ASSETHANDLER :: initialization")
	avatarAssetsHandler := assets.NewAvatarHandler(st)
	if err = fetchAssets(avatarAssetsHandler); err != nil {
		logger.Fatal("ASSETHANDLER :: failed fetching assets: %s", err.Error())
	}

	var ms *mailserver.MailServer
	if cfg.MailServer != nil {
		logger.Info("MAILSERVER :: initialization")
		ms, err = mailserver.NewMailServer(cfg.MailServer, "noreply@myrunes.com", "myrunes")
		if err != nil {
			logger.Fatal("MAILSERVER :: failed connecting to mail account: %s", err.Error())
		}
		logger.Info("MAILSERVER :: started")
	} else {
		logger.Warning("MAILSERVER :: mail server is disabled due to missing configuration")
	}

	var cache caching.CacheMiddleware
	if cfg.Redis != nil && cfg.Redis.Enabled {
		cache = caching.NewRedis(cfg.Redis)
	} else {
		cache = caching.NewInternal()
	}
	cache.SetDatabase(db)

	logger.Info("WEBSERVER :: initialization")
	ws, err := webserver.NewWebServer(db, cache, ms, avatarAssetsHandler, cfg.WebServer)
	if err != nil {
		logger.Fatal("WEBSERVER :: failed creating web server: %s", err.Error())
	}
	go func() {
		if err := ws.ListenAndServeBlocking(); err != nil {
			logger.Fatal("WEBSERVER :: failed starting web server: %s", err.Error())
		}
	}()
	logger.Info("WEBSERVER :: started")

	lct := lifecycletimer.New(24 * time.Second).
		Handle(func() { refetch(avatarAssetsHandler) }).
		Handle(func() { cleanupExpiredRefreshTokens(db) }).
		Start()
	defer lct.Stop()
	logger.Info("LIFECYCLETIMER :: started")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
