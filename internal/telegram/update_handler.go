package telegram

import (
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/text"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// processIncomingUpdates processes incoming updates from the user
func (b *Bot) processIncomingUpdates(update tgbotapi.Update) {
	if update.Message != nil && update.Message.IsCommand() { //When user sends command
		b.handleCommand(update.Message, update.SentFrom().FirstName)
	} else if update.Message != nil && update.Message.Location != nil { //When user sends location
		b.handleLocation(update.Message)
	} else if update.CallbackQuery != nil { //When user choose forecast type or the "repeat last" command
		if update.CallbackQuery.Data != text.CallbackLast {
			b.handleCallbackQuery(update.CallbackQuery)
		} else {
			b.handleCallbackLast(update.CallbackQuery, update.SentFrom().FirstName)
		}
	} else if update.Message != nil && !update.Message.IsCommand() { //When user sends cityname
		b.handleText(update.Message)
	}
}

// handleCommand processes commands from the user.
func (b *Bot) handleCommand(message *tgbotapi.Message, fname string) {
	if message.Text == text.CommandMetricUnits || message.Text == text.CommandNonMetricUnits {
		b.handleUnitsCommand(message)
	} else if message.Text == text.CommandStart {
		b.handleStartCommand(message, fname)
	} else if message.Text == text.CommandHelp {
		b.SendMessage(message.Chat.ID, text.MsgHelp)
	}
}

// handleStartCommand handles the /start command.
func (b *Bot) handleStartCommand(message *tgbotapi.Message, fname string) {
	b.weatherService.CreateUserById(message.Chat.ID)
	b.SendMessage(message.Chat.ID, fmt.Sprintf(text.MsgWelcome, fname)+text.MsgHelp)
}

// handleHelpCommand handles the /help command.
func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	b.SendMessage(message.Chat.ID, text.MsgHelp)
}

// handleUnitsCommand handles the /metric and /non-metric commands.
func (b *Bot) handleUnitsCommand(message *tgbotapi.Message) {
	err := b.weatherService.UserData.SetUserMeasurementSystem(message.Chat.ID, message.Text)
	if err != nil {
		b.SendMessage(message.Chat.ID, text.MsgSetUsersSystemError)
	}
	b.SendMessage(message.Chat.ID, text.MsgMetricUnitChanged)
}

// handleText processes text from the user.
func (b *Bot) handleText(message *tgbotapi.Message) {
	if !containsEmoji(message.Text) {
		err := b.weatherService.UserData.SetUserLastInputCity(message.Chat.ID, message.Text)
		if err != nil {
			b.SendMessage(message.Chat.ID, text.MsgSetUsersCityError)
		}
		b.SendMessageWithInlineKeyboard(message.Chat.ID, text.MsgChooseOption, text.CallbackCurrent, text.CallbackForecast)
	} else {
		b.SendMessage(message.Chat.ID, text.MsgUnsupportedMessageType)
	}
}

// containsEmoji returns true if the text contains emojis.
func containsEmoji(text string) bool {
	for _, char := range text {
		if char >= '\U0001F600' && char <= '\U0001F64F' {
			return true
		}
	}
	return false
}

// handleLocationMessage processes location messages from the user.
func (b *Bot) handleLocation(message *tgbotapi.Message) {
	uLat, uLon := fmt.Sprintf("%f", message.Location.Latitude), fmt.Sprintf("%f", message.Location.Longitude)
	err := b.weatherService.UserData.SetUserLastInputLocation(message.Chat.ID, uLat, uLon)
	if err != nil {
		b.SendMessage(message.Chat.ID, text.MsgSetUsersLocationError)
	}
	b.SendLocationOptions(message.Chat.ID, uLat, uLon)
}

// handleCallbackQuery handles callback queries from the user.
func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {

	_ = b.weatherService.UserData.SetUserLastWeatherCommand(callback.Message.Chat.ID, callback.Data)
	user, err := b.weatherService.UserData.GetUserById(callback.Message.Chat.ID)
	if err != nil {
		b.SendMessage(callback.Message.Chat.ID, err.Error())
	} else {
		b.sendWeather(callback.Message.Chat.ID, user)
	}
}

// handleCallbackQuery handles callback queries from the user with the "repeat last" command.
func (b *Bot) handleCallbackLast(callback *tgbotapi.CallbackQuery, fname string) {

	user, err := b.weatherService.UserData.GetUserById(callback.Message.Chat.ID)
	if err != nil {
		b.SendMessage(callback.Message.Chat.ID, err.Error())
	} else {
		if user.Last == "" {
			b.SendMessage(callback.Message.Chat.ID, fmt.Sprintf(text.MsgLastDataUnavailable, fname))
		} else {
			b.sendWeather(callback.Message.Chat.ID, user)
		}
	}

}

// sendWeather retrieves and sends weather information to the user.
func (b *Bot) sendWeather(chatId int64, user model.UserData) {

	userMessage, err := b.weatherService.WeatherApi.GetWeatherForecast(user)
	if err != nil {
		b.SendMessage(chatId, err.Error())
	} else {
		b.SendMessageWithInlineKeyboard(chatId, userMessage, text.CallbackLast)
	}
}
