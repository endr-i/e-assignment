package conf

import (
	"assignment/logger"
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
	Log    logger.LogConfig
	Port   string `default:"8080"`
}

func GetConfig() Config {
	once.Do(func() {
		configor.New(&configor.Config{
			ENVPrefix: "E",
		}).Load(&config, "./config.yml")
	})
	return config
}
