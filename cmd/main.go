package main

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/logger"
	"SimpleWeatherTgBot/repository"
	"SimpleWeatherTgBot/telegram"
	"SimpleWeatherTgBot/weather_service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.NewConfig()
	repo := repository.NewRepository()
	log := logger.SetupLogger()
	weatherService := weather_service.NewWClient(repo, cfg, log)

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	tBot := telegram.NewBot(botApi, log, weatherService, cfg)
	tBot.Run()
}
