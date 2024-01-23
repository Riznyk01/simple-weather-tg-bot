package telegram

import (
	"SimpleWeatherTgBot/repository"
	"SimpleWeatherTgBot/types"
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
	case types.CommandMetricUnits:
		err := b.weatherService.WeatherControl.SetSystem(message.Chat.ID, true)
		if err != nil {
			b.SendMessage(message.Chat.ID, types.MessageSetUsersSystemError)
			b.log.Errorf("%s: %s", fc, types.MessageSetUsersSystemError)
		}
		b.SendMessage(message.Chat.ID, types.MessageMetricUnitOn)
	case types.CommandNonMetricUnits:
		err := b.weatherService.WeatherControl.SetSystem(message.Chat.ID, false)
		if err != nil {
			b.SendMessage(message.Chat.ID, types.MessageSetUsersSystemError)
			b.log.Errorf("%s: %s", fc, types.MessageSetUsersSystemError)
		}
		b.SendMessage(message.Chat.ID, types.MessageMetricUnitOff)
	case types.CommandStart:
		b.SendMessage(message.Chat.ID, fmt.Sprintf(types.MessageWelcome, fname)+types.MessageHelp)
	case types.CommandHelp:
		b.SendMessage(message.Chat.ID, types.MessageHelp)
	}
}

// handleText processes text from users.
func (b *Bot) handleText(message *tgbotapi.Message) {
	fc := "handleText"

	err := b.weatherService.WeatherControl.SetCity(message.Chat.ID, message.Text)
	if err != nil {
		b.SendMessage(message.Chat.ID, types.MessageSetUsersCityError)
		b.log.Errorf("%s: %s", fc, types.MessageSetUsersCityError)
	}
	err = b.SendMessageWithInlineKeyboard(message.Chat.ID, types.MessageChooseOption, types.CallbackCurrent, types.CallbackForecast)
	if err != nil {
		b.log.Errorf("%s: %v", fc, err)
	}
}

// handleLocationMessage processes location messages from users.
func (b *Bot) handleLocation(message *tgbotapi.Message) {
	fc := "handleLocation"

	uLat, uLon := fmt.Sprintf("%f", message.Location.Latitude), fmt.Sprintf("%f", message.Location.Longitude)
	err := b.weatherService.WeatherControl.SetLocation(message.Chat.ID, uLat, uLon)
	if err != nil {
		b.log.Errorf("%s: %s", fc, types.MessageSetUsersLocationError)
		b.SendMessage(message.Chat.ID, types.MessageSetUsersLocationError)
	}
	err = b.SendLocationOptions(message.Chat.ID, uLat, uLon)
	if err != nil {
		b.log.Errorf("%s: %v", fc, err)
	}
}

// handleCallbackQuery processes callback queries from users.
func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery, fname string) {
	var userMessage string
	var err error
	if callback.Data == types.CallbackLast {
		userMessage, err = b.weatherService.WeatherControl.GetLast(callback.Message.Chat.ID)
	} else {
		userMessage, err = b.weatherService.WeatherControl.SetLast(callback.Message.Chat.ID, callback.Data)
	}
	if err != nil {
		if errors.Is(err, repository.ErrItemIsEmpty) {
			b.SendMessage(callback.Message.Chat.ID, fmt.Sprintf(types.MessageLastDataUnavailable, fname))
		} else {
			b.log.Error(err)
			b.SendMessage(callback.Message.Chat.ID, err.Error())
		}
	} else {
		err = b.SendMessageWithInlineKeyboard(callback.Message.Chat.ID, userMessage, types.CallbackLast)
		if err != nil {
			b.log.Error(err)
		}
	}
}
