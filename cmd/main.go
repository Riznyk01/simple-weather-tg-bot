package main

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/http_client"
	"SimpleWeatherTgBot/internal/logger"
	"SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/service"
	"SimpleWeatherTgBot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logger.SetupLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to get bot config: %s", err.Error())
	}

	httpClient := &http_client.DefaultHTTPClient{
		Timeout: time.Second * 10,
	}

	postgresCfg, err := config.NewConfigPostgres()
	if err != nil {
		log.Fatalf("failed to get DB config: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(postgresCfg)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(log, db)
	weatherService := service.NewService(repo, cfg, log, httpClient)

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	tBot := telegram.NewBot(botApi, log, weatherService, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Info("Bot is starting.")
		tBot.Run()
	}()

	<-stop

	log.Info("Received shutdown signal. Initiating graceful shutdown...")
	tBot.Stop()
	err = db.Close()
	if err != nil {
		log.Error(err)
	}
}
