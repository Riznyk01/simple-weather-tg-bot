package main

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/http_client"
	"SimpleWeatherTgBot/internal/logger"
	"SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/service"
	"SimpleWeatherTgBot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func main() {

	cfg := config.NewConfig()
	mem := repository.NewMemoryStorage()
	repo := repository.NewRepository(mem)
	log := logger.SetupLogger()
	httpClient := &http_client.DefaultHTTPClient{
		Timeout: time.Second * 10,
	}
	weatherService := service.NewService(repo, cfg, log, httpClient)

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	tBot := telegram.NewBot(botApi, log, weatherService, cfg)
	tBot.Run()
}
