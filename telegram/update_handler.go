package telegram

import (
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/weather"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Processes text messages and commands from users.
func (b *Bot) handleUpdateMessage(update tgbotapi.Update) {
	switch {
	case update.Message.Text == types.CommandMetricUnits:
		b.storage.SetSystem(update.Message.Chat.ID, true)
		b.SendMessage(update.Message.Chat.ID, types.MetrikUnitOn)
	case update.Message.Text == types.CommandNonMetricUnits:
		b.storage.SetSystem(update.Message.Chat.ID, false)
		b.SendMessage(update.Message.Chat.ID, types.MetrikUnitOff)
	case update.Message.Text == types.CommandStart:
		n := update.SentFrom()
		greet := fmt.Sprintf("%s%s%s%s", types.WelcomeMessage, n.FirstName, types.WelcomeMessageEnd, types.HelpMessage)
		b.SendMessage(update.Message.Chat.ID, greet)
	case update.Message.Text == types.CommandHelp:
		b.SendMessage(update.Message.Chat.ID, types.HelpMessage)
	default:
		b.storage.SetCity(update.Message.Chat.ID, update.Message.Text)
		err := b.SendMessageWithInlineKeyboard(update.Message.Chat.ID, types.ChooseOptionMessage, types.CommandCurrent, types.CommandForecast)
		if err != nil {
			b.log.Error(err)
		}
	}
}

// Processes location messages from users.
func (b *Bot) handleLocationMessage(update tgbotapi.Update) {
	uLat, uLon := fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
	b.storage.SetLocation(update.Message.Chat.ID, uLat, uLon)
	err := b.SendLocationOptions(update.Message.Chat.ID, uLat, uLon)
	if err != nil {
		b.log.Error(err)
	}
}

// Processes callback queries from users.
func (b *Bot) handleCallbackQuery(update tgbotapi.Update) {
	switch {
	case update.CallbackQuery.Data == types.CommandCurrent || update.CallbackQuery.Data == types.CommandForecast:
		if b.storage.GetCity(update.CallbackQuery.Message.Chat.ID) == "" {
			userMessage := types.MissingCityMessage
			b.SendMessage(update.CallbackQuery.Message.Chat.ID, userMessage)
		} else {
			weatherUrl, err := weather.WeatherUrlByCity(b.storage.GetCity(update.CallbackQuery.Message.Chat.ID), b.cfg.WToken, update.CallbackQuery.Data, b.storage.GetSystem(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				b.log.Error(err)
			}
			b.storage.SetLast(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			userMessage, err := weather.GetWeather(weatherUrl, update.CallbackQuery.Data, b.storage.GetSystem(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				b.log.Error(err)
				userMessage = err.Error()
				b.SendMessage(update.CallbackQuery.Message.Chat.ID, userMessage)
			} else {
				err = b.SendMessageWithInlineKeyboard(update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
				if err != nil {
					b.log.Error(err)
				}
			}
		}
	case update.CallbackQuery.Data == types.CommandForecastLocation || update.CallbackQuery.Data == types.CommandCurrentLocation:
		if b.storage.GetLat(update.CallbackQuery.Message.Chat.ID) == "" && b.storage.GetLon(update.CallbackQuery.Message.Chat.ID) == "" {
			userMessage := types.NoLocationProvidedMessage
			b.SendMessage(update.CallbackQuery.Message.Chat.ID, userMessage)
		} else {
			weatherUrl, err := weather.WeatherUrlByLocation(b.storage.GetLat(update.CallbackQuery.Message.Chat.ID), b.storage.GetLon(update.CallbackQuery.Message.Chat.ID), b.cfg.WToken, update.CallbackQuery.Data, b.storage.GetSystem(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				b.log.Error(err)
			}
			b.storage.SetLast(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			userMessage, err := weather.GetWeather(weatherUrl, update.CallbackQuery.Data, b.storage.GetSystem(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				b.log.Error(err)
			} else {
				err = b.SendMessageWithInlineKeyboard(update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
				if err != nil {
					b.log.Error(err)
				}
			}
		}
	}
}

// Processes the "repeat last" callback query, sends the last weather data.
func (b *Bot) handleLast(update tgbotapi.Update) {
	switch {
	// If the users last requested weather type is empty due to a bot restart.
	case !b.storage.Exists(update.CallbackQuery.Message.Chat.ID):
		name := update.SentFrom()
		b.SendMessage(update.CallbackQuery.Message.Chat.ID, types.LastDataUnavailable+name.FirstName+types.LastDataUnavailableEnd)
	case b.storage.GetLast(update.CallbackQuery.Message.Chat.ID) == types.CommandCurrent || b.storage.GetLast(update.CallbackQuery.Message.Chat.ID) == types.CommandForecast:
		weatherUrl, err := weather.WeatherUrlByCity(b.storage.GetCity(update.CallbackQuery.Message.Chat.ID), b.cfg.WToken, b.storage.GetLast(update.CallbackQuery.Message.Chat.ID), b.storage.GetSystem(update.CallbackQuery.Message.Chat.ID))
		if err != nil {
			b.log.Error(err)
		}
		userMessage, err := weather.GetWeather(weatherUrl, b.storage.GetLast(update.CallbackQuery.Message.Chat.ID), b.storage.GetSystem(update.CallbackQuery.Message.Chat.ID))
		if err != nil {
			b.log.Error(err)
			userMessage = err.Error()
		}
		err = b.SendMessageWithInlineKeyboard(update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
		if err != nil {
			b.log.Error(err)
		}
	case b.storage.GetLast(update.CallbackQuery.Message.Chat.ID) == types.CommandForecastLocation || b.storage.GetLast(update.CallbackQuery.Message.Chat.ID) == types.CommandCurrentLocation:
		weatherUrl, err := weather.WeatherUrlByLocation(b.storage.GetLat(update.CallbackQuery.Message.Chat.ID), b.storage.GetLon(update.CallbackQuery.Message.Chat.ID), b.cfg.WToken, b.storage.GetLast(update.CallbackQuery.Message.Chat.ID), b.storage.GetSystem(update.CallbackQuery.Message.Chat.ID))
		if err != nil {
			b.log.Error(err)
		}
		userMessage, err := weather.GetWeather(weatherUrl, b.storage.GetLast(update.CallbackQuery.Message.Chat.ID), b.storage.GetSystem(update.CallbackQuery.Message.Chat.ID))
		if err != nil {
			b.log.Error(err)
		}
		err = b.SendMessageWithInlineKeyboard(update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
		if err != nil {
			b.log.Error(err)
		}
	}
}
