package telegram

import (
	"SimpleWeatherTgBot/admin"
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/weather_service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot            *tgbotapi.BotAPI
	log            *logrus.Logger
	weatherService *weather_service.WeatherService
	adminService   *admin.AdminService
	cfg            *config.Config
}

func NewBot(bot *tgbotapi.BotAPI, log *logrus.Logger, weatherService *weather_service.WeatherService,
	cfg *config.Config, adminService *admin.AdminService) *Bot {
	return &Bot{
		bot:            bot,
		log:            log,
		weatherService: weatherService,
		adminService:   adminService,
		cfg:            cfg,
	}
}
func (b *Bot) Run() error {
	b.bot.Debug = b.cfg.BotDebug
	b.log.Infof("Authorized on account %s", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		b.processIncomingUpdates(update)
	}
	return nil
}
