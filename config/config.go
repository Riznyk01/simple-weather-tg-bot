package config

import "os"

type Config struct {
	BotToken string
	WToken   string
	LogLevel string //LOG_LEVEL
	LogType  string //TYPE_OF_LOG
}

func NewConfig() *Config {
	return &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
		WToken:   os.Getenv("WEATHER_TOKEN"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		LogType:  os.Getenv("TYPE_OF_LOG"),
	}
}
