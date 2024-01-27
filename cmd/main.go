package main

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/logger"
	"SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/telegram"
	"SimpleWeatherTgBot/internal/user_management_service"
	"SimpleWeatherTgBot/internal/weather_service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	cfg := config.NewConfig()
	mem := repository.NewMemoryStorage()
	repo := repository.NewRepository(mem)
	log := logger.SetupLogger()
	weatherService := weather_service.NewWeatherService(repo, cfg, log)
	userService := user_management_service.NewUserService(repo, cfg, log)

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	tBot := telegram.NewBot(botApi, log, weatherService, cfg, userService)
	tBot.Run()
}
