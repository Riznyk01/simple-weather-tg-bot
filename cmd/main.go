package main

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/http_client"
	"SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/service"
	"SimpleWeatherTgBot/internal/telegram"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var log logr.Logger

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error(err, "Failed to get bot config:")
	}

	httpClient := &http_client.DefaultHTTPClient{
		Timeout: time.Second * 10,
	}

	postgresCfg, err := config.NewConfigPostgres()
	if err != nil {
		log.Error(err, "Failed to get DB config:")
	}

	db, err := repository.NewPostgresDB(postgresCfg)
	if err != nil {
		log.Error(err, "Failed to initialize db:")
	}

	repo := repository.NewRepository(&log, db)
	weatherService := service.NewService(repo, cfg, &log, httpClient)

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Error(err, "Error occurred while creating a new BotAPI instance:")
	}

	tBot := telegram.NewBot(botApi, &log, weatherService, cfg)

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
		log.Error(err, "error occurred while closing the DB")
	}
}
