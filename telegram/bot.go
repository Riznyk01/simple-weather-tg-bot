package telegram

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/storage"
	"SimpleWeatherTgBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	log     *logrus.Logger
	storage *storage.MemoryStorage
	cfg     *config.Config
}

func NewBot(bot *tgbotapi.BotAPI, log *logrus.Logger, storage *storage.MemoryStorage, cfg *config.Config) *Bot {
	return &Bot{
		bot:     bot,
		log:     log,
		storage: storage,
		cfg:     cfg,
	}
}
func (b *Bot) Run() error {
	b.bot.Debug = b.cfg.BotDebug
	b.log.Infof("Authorized on account %s", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		switch {
		case update.Message != nil && update.Message.Location == nil:
			//When user sends command or cityname
			b.handleUpdateMessage(update)
		case update.Message != nil && update.Message.Location != nil:
			//When user sends location
			b.handleLocationMessage(update)
		case update.Message == nil && update.CallbackQuery != nil:
			if update.CallbackQuery.Data != types.CommandLast {
				//When user choose forecast type
				b.handleCallbackQuery(update)
			} else {
				//When user choose last forecast
				b.handleLast(update)
			}
		}
	}
	return nil
}
