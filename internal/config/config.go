package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/zekroTJA/lol-runes/internal/webserver"
	"github.com/zekrotja/lol-runes/internal/database"
)

type Main struct {
	MongoDB   *database.MongoConfig `json:"mongodb"`
	WebServer *webserver.Config     `json:"webserver"`
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
			Username: "kirby",
			AuthDB:   "kirby",
			DataDB:   "kirby",
		},
		WebServer: &webserver.Config{
			Addr: ":443",
			TLS: &webserver.TLSConfig{
				Enabled: true,
			},
		},
	}

	data, err := yaml.Marshal(def)

	basePath := path.Base(loc)
	if _, err = os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0750)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(loc, data, 0750)
	return err
}
