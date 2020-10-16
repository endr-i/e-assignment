package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type LogConfig struct {
	File  string `default:""`
	Level int    `default:"0"` // 0 - 6
}

func InitLogger(config LogConfig) {
	log.SetFormatter(&log.JSONFormatter{})
	if config.File == "" {
		log.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Error("cannot open logger file")
			log.SetOutput(os.Stdout)
		}
		log.SetOutput(file)
	}
	log.SetLevel(log.Level(config.Level))
}
