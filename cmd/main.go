package main

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/logger"
	repository2 "SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/telegram"
	"SimpleWeatherTgBot/internal/weather_service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	cfg := config.NewConfig()
	mem := repository2.NewMemoryStorage()
	repo := repository2.NewRepository(mem)
	log := logger.SetupLogger()
	weatherService := weather_service.NewWClient(repo, cfg, log)

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	tBot := telegram.NewBot(botApi, log, weatherService, cfg)
	tBot.Run()
}
