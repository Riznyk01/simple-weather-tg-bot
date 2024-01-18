package main

import (
	"SimpleWeatherTgBot/admin"
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/logger"
	"SimpleWeatherTgBot/repository"
	"SimpleWeatherTgBot/telegram"
	"SimpleWeatherTgBot/weather_service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	cfg := config.NewConfig()
	mem := repository.NewMemoryStorage()
	repo := repository.NewRepository(mem)
	log := logger.SetupLogger()
	adm := admin.NewAdminService(repo, log, cfg)
	weatherService := weather_service.NewWClient(repo, cfg, log)

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	tBot := telegram.NewBot(botApi, log, weatherService, cfg, adm)
	tBot.Run()
}
