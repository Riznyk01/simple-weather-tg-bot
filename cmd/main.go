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

	cfg := config.NewConfig()
	log := logger.SetupLogger()
	httpClient := &http_client.DefaultHTTPClient{
		Timeout: time.Second * 10,
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_Host"),
		Port:     os.Getenv("DB_Port"),
		Username: os.Getenv("DB_Username"),
		Password: os.Getenv("DB_Password"),
		DBName:   os.Getenv("DB_Name"),
		SSLMode:  os.Getenv("DB_SSLMode"),
	})
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
