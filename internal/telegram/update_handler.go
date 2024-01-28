package telegram

import (
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/repository"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) processIncomingUpdates(update tgbotapi.Update) {
	b.log.Debug(update)
	switch {
	case update.Message != nil && update.Message.IsCommand(): //When user sends command
		b.handleCommand(update.Message, update.SentFrom().FirstName)
	case update.Message != nil && update.Message.Location != nil: //When user sends location
		b.handleLocation(update.Message)
	case update.CallbackQuery != nil: //When user choose forecast type or the "repeat last" command
		b.handleCallbackQuery(update.CallbackQuery, update.SentFrom().FirstName)
	default: //When user sends cityname
		b.handleText(update.Message)
	}
}

// handleCommand processes command from users.
func (b *Bot) handleCommand(message *tgbotapi.Message, fname string) {
	fc := "handleCommand"

	switch message.Text {
	case model.CommandMetricUnits, model.CommandNonMetricUnits:
		err := b.userService.SetSystem(message.Chat.ID, message.Text)
		if err != nil {
			b.SendMessage(message.Chat.ID, model.MessageSetUsersSystemError)
			b.log.Errorf("%s: %s", fc, model.MessageSetUsersSystemError)
		}
		b.SendMessage(message.Chat.ID, model.MessageMetricUnitChanged)
	case model.CommandStart:
		b.SendMessage(message.Chat.ID, fmt.Sprintf(model.MessageWelcome, fname)+model.MessageHelp)
	case model.CommandHelp:
		b.SendMessage(message.Chat.ID, model.MessageHelp)
	}
}

// handleText processes text from users.
func (b *Bot) handleText(message *tgbotapi.Message) {
	fc := "handleText"

	err := b.userService.SetCity(message.Chat.ID, message.Text)
	if err != nil {
		b.SendMessage(message.Chat.ID, model.MessageSetUsersCityError)
		b.log.Errorf("%s: %s", fc, model.MessageSetUsersCityError)
	}
	b.SendMessageWithInlineKeyboard(message.Chat.ID, model.MessageChooseOption, model.CallbackCurrent, model.CallbackForecast)
}

// handleLocationMessage processes location messages from users.
func (b *Bot) handleLocation(message *tgbotapi.Message) {
	uLat, uLon := fmt.Sprintf("%f", message.Location.Latitude), fmt.Sprintf("%f", message.Location.Longitude)
	err := b.userService.SetLocation(message.Chat.ID, uLat, uLon)
	if err != nil {
		b.SendMessage(message.Chat.ID, model.MessageSetUsersLocationError)
	}
	b.SendLocationOptions(message.Chat.ID, uLat, uLon)
}

// handleCallbackQuery processes callback queries from users.
func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery, fname string) {
	fc := "handleCallbackQuery"
	var last bool
	b.log.Debugf("%s: callback: %s", fc, callback.Data)
	userMessage, err := b.weatherService.WeatherControl.GetWeatherForecast(callback.Message.Chat.ID, callback.Data)
	if err != nil {
		last = true
		if errors.Is(err, repository.ErrItemIsEmpty) {
			b.SendMessage(callback.Message.Chat.ID, fmt.Sprintf(model.MessageLastDataUnavailable, fname))
		} else {
			b.SendMessage(callback.Message.Chat.ID, err.Error())
		}
	}
	if !last {
		b.SendMessageWithInlineKeyboard(callback.Message.Chat.ID, userMessage, model.CallbackLast)
	}
}
