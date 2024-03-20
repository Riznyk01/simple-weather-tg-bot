package telegram

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/service"
	"SimpleWeatherTgBot/internal/text"
	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	botApi  *tgbotapi.BotAPI
	cfg     *config.Config
	log     *logr.Logger
	service *service.Service
}

func NewBot(botApi *tgbotapi.BotAPI, log *logr.Logger, cfg *config.Config, service *service.Service) *Bot {
	return &Bot{
		botApi:  botApi,
		cfg:     cfg,
		log:     log,
		service: service,
	}
}
func (b *Bot) Run() {
	b.botApi.Debug = b.cfg.BotDebug
	b.log.Info("Authorized on account", "account", b.botApi.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botApi.GetUpdatesChan(u)

	for update := range updates {
		go b.processIncomingUpdates(update)
	}
}

// processIncomingUpdates processes incoming updates from the user
func (b *Bot) processIncomingUpdates(update tgbotapi.Update) {

	var m model.UserMessage
	var err error

	if update.Message != nil && update.Message.IsCommand() { // When user sends command

		m, _ = b.service.Handler.HandleCommand(update.Message, update.SentFrom().FirstName)
		b.sendMessage(update.Message.Chat.ID, m.Text)
	} else if update.Message != nil && update.Message.Location != nil { // When user sends location
		m, err = b.service.Handler.HandleLocation(update.Message)
		if err != nil {
			b.sendMessage(update.Message.Chat.ID, m.Text)
		} else {
			b.sendMessageWithInlineKeyboard(update.Message.Chat.ID, m)
		}
	} else if update.CallbackQuery != nil { // When user choose forecast type or the "repeat last" command
		if update.CallbackQuery.Data != text.CallbackLast {
			m, err = b.service.Handler.HandleCallbackQuery(update.CallbackQuery)
			if m.Buttons == nil {
				b.sendMessage(update.CallbackQuery.Message.Chat.ID, text.TryAnother)
			} else {
				b.sendMessageWithInlineKeyboard(update.CallbackQuery.Message.Chat.ID, m)
			}
		} else {
			m, err = b.service.Handler.HandleCallbackLast(update.CallbackQuery, update.SentFrom().FirstName)
			b.sendMessageWithInlineKeyboard(update.CallbackQuery.Message.Chat.ID, m)
		}
	} else if update.Message != nil && !update.Message.IsCommand() { // When user sends cityname
		m, err = b.service.Handler.HandleText(update.Message)
		if m.Buttons != nil {
			b.sendMessageWithInlineKeyboard(update.Message.Chat.ID, m)
		} else {
			b.sendMessage(update.Message.Chat.ID, text.MsgSetUsersCityError)
		}
	}
}

// sendMessageWithInlineKeyboard sends a message with text and inline keyboard for service type selection
func (b *Bot) sendMessageWithInlineKeyboard(chatID int64, userMsg model.UserMessage) {

	msg := tgbotapi.NewMessage(chatID, userMsg.Text)

	var inlineButtons []tgbotapi.InlineKeyboardButton
	for _, button := range userMsg.Buttons {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(inlineButtons...))
	msg.ReplyMarkup = inlineKeyboard
	msg.ParseMode = "HTML"
	_, err := b.botApi.Send(msg)
	if err != nil {
		b.log.Error(err, text.ErrWhileSendingInline)
	}
}

func (b *Bot) sendMessage(chatID int64, msgText string) {
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = "HTML"
	_, err := b.botApi.Send(msg)
	if err != nil {
		b.log.Error(err, text.ErrWhileSendingMsg)
	}
}

func (b *Bot) Stop() {
	b.botApi.StopReceivingUpdates()
	b.log.Info("Graceful shutdown complete.")
}
