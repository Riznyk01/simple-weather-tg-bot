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
		b.weatherService.WeatherUserControl.SetSystem(chatId, true)
		b.SendMessage(chatId, types.MetricUnitOn)
	case types.CommandNonMetricUnits:
		b.weatherService.WeatherUserControl.SetSystem(chatId, false)
		b.SendMessage(chatId, types.MetricUnitOff)
	case types.CommandStart:
		n := update.SentFrom()
		greet := fmt.Sprintf("%s%s%s%s", types.WelcomeMessage, n.FirstName, types.WelcomeMessageEnd, types.HelpMessage)
		b.SendMessage(chatId, greet)
	case types.CommandHelp:
		b.SendMessage(chatId, types.HelpMessage)
	default:
		b.weatherService.WeatherUserControl.SetCity(chatId, update.Message.Text)
		err := b.SendMessageWithInlineKeyboard(chatId, types.ChooseOptionMessage, types.CommandCurrent, types.CommandForecast)
		if err != nil {
			b.log.Error(err)
		}
	}
}

// handleLocationMessage processes location messages from users.
func (b *Bot) handleLocationMessage(update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	uLat, uLon := fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
	b.weatherService.WeatherUserControl.SetLocation(chatId, uLat, uLon)
	err := b.SendLocationOptions(chatId, uLat, uLon)
	if err != nil {
		b.log.Error(err)
	}
}

// handleCallbackQuery processes callback queries from users.
func (b *Bot) handleCallbackQuery(update tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID
	weatherCommand := update.CallbackQuery.Data
	userMessage, err := b.weatherService.WeatherUserControl.SetLast(chatId, weatherCommand)
	b.handleCallbackQueryHandlingError(update, userMessage, err)
}

// handleCallbackQueryLast processes the "repeat last" callback query, sends the last weather data.
func (b *Bot) handleCallbackQueryLast(update tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID
	userMessage, err := b.weatherService.WeatherUserControl.GetLast(chatId)
	b.handleCallbackQueryHandlingError(update, userMessage, err)
}

// handleCallbackQueryHandlingError handles errors in callback query processing.
func (b *Bot) handleCallbackQueryHandlingError(update tgbotapi.Update, userMessage string, err error) {
	if userMessage == "empty" {
		n := update.SentFrom()
		b.SendMessage(update.Message.Chat.ID, fmt.Sprintf(types.LastDataUnavailable+n.FirstName+types.LastDataUnavailableEnd))
	} else if err != nil {
		b.log.Error(err)
		b.SendMessage(update.Message.Chat.ID, err.Error())
	} else {
		err = b.SendMessageWithInlineKeyboard(update.Message.Chat.ID, userMessage, types.CommandLast)
		if err != nil {
			b.log.Error(err)
		}
	}
}
