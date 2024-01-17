package telegram

import (
	"SimpleWeatherTgBot/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleUpdateMessage processes text messages and commands from users.
func (b *Bot) handleUpdateMessage(update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	switch update.Message.Text {
	case types.CommandMetricUnits:
		err := b.weatherService.WeatherUserControl.SetSystem(chatId, true)
		if err != nil {
			b.SendMessage(update.Message.Chat.ID, types.SetUsersSystemError)
			b.log.Error(types.SetUsersSystemError)
		}
		b.SendMessage(chatId, types.MetricUnitOn)
	case types.CommandNonMetricUnits:
		err := b.weatherService.WeatherUserControl.SetSystem(chatId, false)
		if err != nil {
			b.SendMessage(update.Message.Chat.ID, types.SetUsersSystemError)
			b.log.Error(types.SetUsersSystemError)
		}
		b.SendMessage(chatId, types.MetricUnitOff)
	case types.CommandStart:
		n := update.SentFrom()
		greet := fmt.Sprintf(types.WelcomeMessage, n.FirstName) + types.HelpMessage
		b.SendMessage(chatId, greet)
	case types.CommandHelp:
		b.SendMessage(chatId, types.HelpMessage)
	default:
		err := b.weatherService.WeatherUserControl.SetCity(chatId, update.Message.Text)
		if err != nil {
			b.SendMessage(update.Message.Chat.ID, types.SetUsersCityError)
			b.log.Error(types.SetUsersCityError)
		}
		err = b.SendMessageWithInlineKeyboard(chatId, types.ChooseOptionMessage, types.CommandCurrent, types.CommandForecast)
		if err != nil {
			b.log.Error(err)
		}
	}
}

// handleLocationMessage processes location messages from users.
func (b *Bot) handleLocationMessage(update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	uLat, uLon := fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
	err := b.weatherService.WeatherUserControl.SetLocation(chatId, uLat, uLon)
	if err != nil {
		b.log.Error(types.SetUsersLocationError)
		b.SendMessage(update.Message.Chat.ID, types.SetUsersLocationError)
	}
	err = b.SendLocationOptions(chatId, uLat, uLon)
	if err != nil {
		b.log.Error(err)
	}
}

// handleCallbackQuery processes callback queries from users.
func (b *Bot) handleCallbackQuery(update tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID
	weatherCommand := update.CallbackQuery.Data
	userMessage, err := b.weatherService.WeatherUserControl.SetLast(chatId, weatherCommand)
	b.handleCallbackQueryHandlingError(update.SentFrom().FirstName, userMessage, chatId, err)
}

// handleCallbackQueryLast processes the "repeat last" callback query, sends the last weather data.
func (b *Bot) handleCallbackQueryLast(update tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID
	userMessage, err := b.weatherService.WeatherUserControl.GetLast(chatId)
	b.handleCallbackQueryHandlingError(update.SentFrom().FirstName, userMessage, chatId, err)
}

// handleCallbackQueryHandlingError handles errors in callback query processing.
func (b *Bot) handleCallbackQueryHandlingError(name, userMessage string, chatId int64, err error) {
	if userMessage == "empty" {
		b.SendMessage(chatId, fmt.Sprintf(types.LastDataUnavailable, name))
	} else if err != nil {
		b.log.Error(err)
		b.SendMessage(chatId, err.Error())
	} else {
		err = b.SendMessageWithInlineKeyboard(chatId, userMessage, types.CommandLast)
		if err != nil {
			b.log.Error(err)
		}
	}
}
