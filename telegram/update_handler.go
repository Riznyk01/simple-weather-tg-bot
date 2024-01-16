package telegram

import (
	"SimpleWeatherTgBot/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Processes text messages and commands from users.
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
	//}
	//
	//// Processes location messages from users.
	//func (b *Bot) handleLocationMessage(update tgbotapi.Update) {
	//	chatId := update.Message.Chat.ID
	//	uLat, uLon := fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
	//	b.storage.SetLocation(chatId, uLat, uLon)
	//	err := b.SendLocationOptions(chatId, uLat, uLon)
	//	if err != nil {
	//		b.log.Error(err)
	//	}
	//}
	//
	//// Processes callback queries from users.
	//func (b *Bot) handleCallbackQuery(update tgbotapi.Update) {
	//	chatId := update.CallbackQuery.Message.Chat.ID
	//	callback := update.CallbackQuery.Data
	//	param := make(map[string]string)
	//	switch {
	//	case callback == types.CommandCurrent || callback == types.CommandForecast:
	//		param["city"] = b.storage.GetCity(chatId)
	//	case callback == types.CommandForecastLocation || callback == types.CommandCurrentLocation:
	//		param["lat"] = b.storage.GetLat(chatId)
	//		param["lon"] = b.storage.GetLon(chatId)
	//	}
	//	weatherUrl, err := weather.GenerateWeatherUrl(param, b.cfg.WToken, callback, b.storage.GetSystem(chatId))
	//	if err != nil {
	//		b.log.Error(err)
	//	}
	//	b.storage.SetLast(chatId, callback)
	//	userMessage, err := weather.GetWeather(weatherUrl, callback, b.storage.GetSystem(chatId))
	//	if err != nil {
	//		b.log.Error(err)
	//		b.SendMessage(chatId, err.Error())
	//	} else {
	//		err = b.SendMessageWithInlineKeyboard(chatId, userMessage, types.CommandLast)
	//		if err != nil {
	//			b.log.Error(err)
	//		}
	//	}
	//}
	//
	//// Processes the "repeat last" callback query, sends the last weather data.
	//func (b *Bot) handleLast(update tgbotapi.Update) {
	//	chatId := update.CallbackQuery.Message.Chat.ID
	//	param := make(map[string]string)
	//	// If the users last requested weather type is empty due to a bot restart.
	//	if !b.storage.Exists(chatId) {
	//		name := update.SentFrom()
	//		b.SendMessage(chatId, types.LastDataUnavailable+name.FirstName+types.LastDataUnavailableEnd)
	//	} else {
	//		switch {
	//		case b.storage.GetLast(chatId) == types.CommandCurrent || b.storage.GetLast(chatId) == types.CommandForecast:
	//			param["city"] = b.storage.GetCity(chatId)
	//		case b.storage.GetLast(chatId) == types.CommandForecastLocation || b.storage.GetLast(chatId) == types.CommandCurrentLocation:
	//			param["lat"] = b.storage.GetLat(chatId)
	//			param["lon"] = b.storage.GetLon(chatId)
	//		}
	//		weatherUrl, err := weather.GenerateWeatherUrl(param, b.cfg.WToken, b.storage.GetLast(chatId), b.storage.GetSystem(chatId))
	//		if err != nil {
	//			b.log.Error(err)
	//		}
	//		userMessage, err := weather.GetWeather(weatherUrl, b.storage.GetLast(chatId), b.storage.GetSystem(chatId))
	//		if err != nil {
	//			b.log.Error(err)
	//			userMessage = err.Error()
	//		}
	//		err = b.SendMessageWithInlineKeyboard(chatId, userMessage, types.CommandLast)
	//		if err != nil {
	//			b.log.Error(err)
	//		}
	//	}
}
