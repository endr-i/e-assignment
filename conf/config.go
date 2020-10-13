package conf

import (
	"assignment/pg"
	"assignment/server"
	"github.com/jinzhu/configor"
	"sync"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	Server server.Config
	DB     pg.Config
}

func GetConfig() Config {
	once.Do(func() {
		configor.New(&configor.Config{
			ENVPrefix: "EXNESS",
		}).Load(&config, "./config.yml")
	})
	return config
}
