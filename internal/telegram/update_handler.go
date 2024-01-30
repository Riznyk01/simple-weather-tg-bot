package telegram

import (
	"SimpleWeatherTgBot/internal/model"
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
		if update.CallbackQuery.Data != model.CallbackLast {
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
	if message.Text == model.CommandMetricUnits || message.Text == model.CommandNonMetricUnits {
		b.handleUnitsCommand(message)
	} else if message.Text == model.CommandStart {
		b.handleStartCommand(message, fname)
	} else if message.Text == model.CommandHelp {
		b.SendMessage(message.Chat.ID, model.MessageHelp)
	}
}

// handleStartCommand handles the /start command.
func (b *Bot) handleStartCommand(message *tgbotapi.Message, fname string) {
	b.weatherService.CreateUser(message.Chat.ID)
	b.SendMessage(message.Chat.ID, fmt.Sprintf(model.MessageWelcome, fname)+model.MessageHelp)
}

// handleHelpCommand handles the /help command.
func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	b.SendMessage(message.Chat.ID, model.MessageHelp)
}

// handleUnitsCommand handles the /metric and /non-metric commands.
func (b *Bot) handleUnitsCommand(message *tgbotapi.Message) {
	err := b.weatherService.UserData.SetSystem(message.Chat.ID, message.Text)
	if err != nil {
		b.SendMessage(message.Chat.ID, model.MessageSetUsersSystemError)
	}
	b.SendMessage(message.Chat.ID, model.MessageMetricUnitChanged)
}

// handleText processes text from the user.
func (b *Bot) handleText(message *tgbotapi.Message) {
	if !containsEmoji(message.Text) {
		err := b.weatherService.UserData.SetCity(message.Chat.ID, message.Text)
		if err != nil {
			b.SendMessage(message.Chat.ID, model.MessageSetUsersCityError)
		}
		b.SendMessageWithInlineKeyboard(message.Chat.ID, model.MessageChooseOption, model.CallbackCurrent, model.CallbackForecast)
	} else {
		b.SendMessage(message.Chat.ID, model.MessageUnsupportedMessageType)
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
	err := b.weatherService.UserData.SetLocation(message.Chat.ID, uLat, uLon)
	if err != nil {
		b.SendMessage(message.Chat.ID, model.MessageSetUsersLocationError)
	}
	b.SendLocationOptions(message.Chat.ID, uLat, uLon)
}

// handleCallbackQuery handles callback queries from the user.
func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {

	_ = b.weatherService.UserData.SetLastWeatherCommand(callback.Message.Chat.ID, callback.Data)
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
			b.SendMessage(callback.Message.Chat.ID, fmt.Sprintf(model.MessageLastDataUnavailable, fname))
		} else {
			b.sendWeather(callback.Message.Chat.ID, user)
		}
	}

}

// sendWeather retrieves and sends weather information to the user.
func (b *Bot) sendWeather(chatId int64, user model.UserData) {
	fc := "sendWeather"
	userMessage, err := b.weatherService.WeatherApi.GetWeatherForecast(user)
	if err != nil {
		b.log.Errorf("%s: %v", fc, err)
		b.SendMessage(chatId, err.Error())
	} else {
		b.SendMessageWithInlineKeyboard(chatId, userMessage, model.CallbackLast)
	}
}
