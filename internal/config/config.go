package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/mailserver"
	"github.com/myrunes/backend/internal/webserver"
)

type Main struct {
	MongoDB    *database.MongoConfig `json:"mongodb"`
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
	fmt.Println(basePath)
	if _, err = os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0750)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(loc, data, 0750)
	return err
}
