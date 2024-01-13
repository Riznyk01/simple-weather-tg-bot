package main

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/logger"
	"SimpleWeatherTgBot/storage"
	"SimpleWeatherTgBot/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log := logger.SetupLogger()
	stor := storage.NewMemoryStorage()
	cfg := config.NewConfig()

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	tBot := telegram.NewBot(botApi, log, stor, cfg)
	tBot.Run()
}
