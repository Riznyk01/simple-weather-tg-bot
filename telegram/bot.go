package telegram

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/storage"
	"SimpleWeatherTgBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"log"
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

	//bot.Debug = debug
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {

		switch {
		case update.Message != nil && update.Message.Location == nil:
			b.handleUpdateMessage(update)
		case update.Message != nil && update.Message.Location != nil:
			b.handleLocationMessage(update)
		case update.Message == nil && update.CallbackQuery != nil:
			if update.CallbackQuery.Data != types.CommandLast {
				b.handleCallbackQuery(update)
			} else {
				b.handleLast(update)
			}

		}
	}
	return nil
}
