package telegram

import (
	"SimpleWeatherTgBot/types"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrLastEmpty = errors.New("the item is empty")

func (b *Bot) processIncomingUpdates(update tgbotapi.Update) {
	switch {
	case update.Message != nil && update.Message.ForwardFrom != nil:
		b.handleReply(update)
	case update.Message != nil && update.Message.Location == nil:
		//When user sends command or cityname
		b.handleTextMessage(update)
	case update.Message != nil && update.Message.Location != nil:
		//When user sends location
		b.handleLocationMessage(update)
	case update.Message == nil && update.CallbackQuery != nil:
		if update.CallbackQuery.Data != types.CommandLast {
			//When user choose forecast type
			b.handleCallbackQuery(update)
		} else {
			//When user choose last forecast
			b.handleCallbackQueryLast(update)
		}
	}
}

// handleTextMessage processes text messages and commands from users.
func (b *Bot) handleTextMessage(update tgbotapi.Update) {
	fc := "handleTextMessage"
	chatId := update.Message.Chat.ID
	b.infoLogger(fc, chatId, update)

	switch update.Message.Text {
	case types.CommandMetricUnits:
		err := b.weatherService.WeatherControl.SetSystem(chatId, true)
		if err != nil {
			b.SendMessage(chatId, types.SetUsersSystemError)
			b.log.Error(types.SetUsersSystemError)
		}
		b.SendMessage(chatId, types.MetricUnitOn)
	case types.CommandNonMetricUnits:
		err := b.weatherService.WeatherControl.SetSystem(chatId, false)
		if err != nil {
			b.SendMessage(chatId, types.SetUsersSystemError)
			b.log.Error(types.SetUsersSystemError)
		}
		b.SendMessage(chatId, types.MetricUnitOff)
	case types.CommandStart:
		n := update.SentFrom()
		greet := fmt.Sprintf(types.WelcomeMessage, n.FirstName) + types.HelpMessage
		b.SendMessage(chatId, greet)
	case types.CommandHelp:
		b.SendMessage(chatId, types.HelpMessage)
	case types.CommandBan:
		//
	case types.CommandUnBan:
		//
	default:
		err := b.weatherService.WeatherControl.SetCity(chatId, update.Message.Text)
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

func (b *Bot) handleReply(update tgbotapi.Update) {
	fc := "handleReply"
	chatId := update.Message.Chat.ID
	b.infoLogger(fc, chatId, update)
	b.weatherService.SetRepliedUserId(update.Message.Chat.ID, update.Message.ForwardFrom.ID)
	replId, _ := b.weatherService.GetRepliedUserId(update.Message.Chat.ID)
	b.SendMessage(chatId, fmt.Sprintf("%v", replId))
}

// handleLocationMessage processes location messages from users.
func (b *Bot) handleLocationMessage(update tgbotapi.Update) {
	fc := "handleLocationMessage"
	chatId := update.Message.Chat.ID
	b.infoLogger(fc, chatId, update)
	uLat, uLon := fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
	err := b.weatherService.WeatherControl.SetLocation(chatId, uLat, uLon)
	if err != nil {
		b.log.Error(types.SetUsersLocationError)
		b.SendMessage(chatId, types.SetUsersLocationError)
	}
	err = b.SendLocationOptions(chatId, uLat, uLon)
	if err != nil {
		b.log.Error(err)
	}
}

// handleCallbackQuery processes callback queries from users.
func (b *Bot) handleCallbackQuery(update tgbotapi.Update) {
	fc := "handleCallbackQuery"
	chatId := update.CallbackQuery.Message.Chat.ID
	b.infoLogger(fc, chatId, update)
	weatherCommand := update.CallbackQuery.Data
	userMessage, err := b.weatherService.WeatherControl.SetLast(chatId, weatherCommand)
	b.handleCallbackQueryHandlingError(update.SentFrom().FirstName, userMessage, chatId, err)
}

// handleCallbackQueryLast processes the "repeat last" callback query, sends the last weather data.
func (b *Bot) handleCallbackQueryLast(update tgbotapi.Update) {
	fc := "handleCallbackQueryLast"
	chatId := update.CallbackQuery.Message.Chat.ID
	b.infoLogger(fc, chatId, update)
	userMessage, err := b.weatherService.WeatherControl.GetLast(chatId)
	b.handleCallbackQueryHandlingError(update.SentFrom().FirstName, userMessage, chatId, err)
}

// handleCallbackQueryHandlingError handles errors in callback query processing.
func (b *Bot) handleCallbackQueryHandlingError(name, userMessage string, chatId int64, err error) {
	if errors.Is(err, ErrLastEmpty) {
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

func (b *Bot) infoLogger(fc string, chatId int64, update tgbotapi.Update) {
	RequestsCount, err := b.weatherService.AddRequestsCount(chatId)
	if err != nil {
		b.log.Error(err)
	}
	var action string
	switch {
	case update.CallbackQuery != nil:
		action = fmt.Sprintf(" callback: %s", update.CallbackQuery.Data)
	case update.Message != nil:
		action = fmt.Sprintf(" message: %s", update.Message.Text)
	case update.Message.Location != nil:
		action = fmt.Sprintf(" location: %v", update.Message.Location)
	}
	b.log.Debug(fc, " ID:", chatId, " ", update.SentFrom().FirstName, update.SentFrom().LastName,
		" @", update.SentFrom().UserName, " req.count: ", RequestsCount, action)
}
