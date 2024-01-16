package main

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/logger"
	"SimpleWeatherTgBot/repository"
	"SimpleWeatherTgBot/telegram"
	"SimpleWeatherTgBot/weather"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	repo := repository.NewRepository()
	weatherService := weather.NewWClient(repo)
	log := logger.SetupLogger()
	//stor := repository.NewMemoryStorage()
	cfg := config.NewConfig()

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	tBot := telegram.NewBot(botApi, log, weatherService, cfg)
	tBot.Run()
}
