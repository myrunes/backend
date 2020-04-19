package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/myrunes/backend/internal/caching"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/mailserver"
	"github.com/myrunes/backend/internal/webserver"
)

type Main struct {
	MongoDB    *database.MongoConfig `json:"mongodb"`
	Redis      *caching.RedisConfig  `json:"redis"`
	WebServer  *webserver.Config     `json:"webserver"`
	MailServer *mailserver.Config    `json:"mailserver"`
}

func Open(loc string) (*Main, error) {
	data, err := ioutil.ReadFile(loc)
	if os.IsNotExist(err) {
		err = cretaeDefault(loc)
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	cfg := new(Main)
	err = yaml.Unmarshal(data, cfg)
	return cfg, err
}

func cretaeDefault(loc string) error {
	def := &Main{
		MongoDB: &database.MongoConfig{
			Host:     "localhost",
			Port:     "27017",
			Username: "lol-runes",
			AuthDB:   "lol-runes",
			DataDB:   "lol-runes",
		},
		Redis: &caching.RedisConfig{
			Enabled: false,
			Addr:    "localhost:6379",
			DB:      0,
		},
		WebServer: &webserver.Config{
			Addr:       ":443",
			PublicAddr: "https://myrunes.com",
			TLS: &webserver.TLSConfig{
				Enabled: true,
			},
		},
		MailServer: &mailserver.Config{
			Port: 465,
		},
	}

	data, err := yaml.Marshal(def)

	basePath := path.Dir(loc)
	if _, err = os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0750)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(loc, data, 0750)
	return err
}
