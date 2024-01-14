package config

import "os"

type Config struct {
	BotToken string
	WToken   string
	LogLevel string
	LogType  string
	BotDebug bool
}

func NewConfig() *Config {

	return &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
		WToken:   os.Getenv("WEATHER_TOKEN"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		LogType:  os.Getenv("TYPE_OF_LOG"),
		BotDebug: DebugStrToBool(os.Getenv("BOT_DEBUG")),
	}
}
func DebugStrToBool(envDebugVar string) bool {
	return envDebugVar == "true"
}
