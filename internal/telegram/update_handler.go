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
	err = b.SendMessageWithInlineKeyboard(message.Chat.ID, model.MessageChooseOption, model.CallbackCurrent, model.CallbackForecast)
	if err != nil {
		b.log.Errorf("%s: %v", fc, err)
	}
}

// handleLocationMessage processes location messages from users.
func (b *Bot) handleLocation(message *tgbotapi.Message) {
	fc := "handleLocation"

	uLat, uLon := fmt.Sprintf("%f", message.Location.Latitude), fmt.Sprintf("%f", message.Location.Longitude)
	err := b.userService.SetLocation(message.Chat.ID, uLat, uLon)
	if err != nil {
		b.log.Errorf("%s: %s", fc, model.MessageSetUsersLocationError)
		b.SendMessage(message.Chat.ID, model.MessageSetUsersLocationError)
	}
	err = b.SendLocationOptions(message.Chat.ID, uLat, uLon)
	if err != nil {
		b.log.Errorf("%s: %v", fc, err)
	}
}

// handleCallbackQuery processes callback queries from users.
func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery, fname string) {
	fc := "handleCallbackQuery"

	command := callback.Data
	var userMessage string
	var err error
	var getLastError bool
	if callback.Data == model.CallbackLast {
		command, err = b.userService.GetLastWeatherCommand(callback.Message.Chat.ID)
		if err != nil {
			if errors.Is(err, repository.ErrItemIsEmpty) {
				b.SendMessage(callback.Message.Chat.ID, fmt.Sprintf(model.MessageLastDataUnavailable, fname))
				getLastError = true
			} else {
				b.log.Errorf("%s: %v", fc, err)
				getLastError = true
				b.SendMessage(callback.Message.Chat.ID, err.Error())
			}
		}
	}
	if !getLastError {
		userMessage, err = b.weatherService.WeatherControl.GetWeatherForecast(callback.Message.Chat.ID, command)
		if err != nil {
			b.log.Error(err)
			b.SendMessage(callback.Message.Chat.ID, err.Error())
		} else {
			err = b.SendMessageWithInlineKeyboard(callback.Message.Chat.ID, userMessage, model.CallbackLast)
			if err != nil {
				b.log.Error(err)
			}
		}
	}
}
