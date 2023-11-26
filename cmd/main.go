package main

import (
	"SimpleWeatherTgBot/lib/e"
	"SimpleWeatherTgBot/logger"
	"SimpleWeatherTgBot/storage"
	"SimpleWeatherTgBot/telegram"
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/weather"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	userMessage, tWeather string
	err                   error
	s                     storage.Storage
)

func main() {

	err := godotenv.Load(".env.dev")
	if err != nil {
		logger.ForError(err)
	}
	t := os.Getenv("BOT_TOKEN")
	tWeather = os.Getenv("WEATHER_KEY")

	bot, err := tgbotapi.NewBotAPI(t)
	if err != nil {
		logger.ForError(e.Wrap("", err))
	}
	bot.Debug = false
	log.SetOutput(os.Stderr)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	s = storage.NewMemoryStorage()

	for update := range updates {
		switch {
		case update.Message != nil && update.Message.Location == nil:
			go handleUpdateMessage(bot, update)
		case update.Message != nil && update.Message.Location != nil:
			go handleLocationMessage(bot, update)
		case update.Message == nil && update.CallbackQuery != nil:
			if update.CallbackQuery.Data != types.CommandLast {
				go handleCallbackQuery(bot, update)
			} else {
				go handleLast(bot, update)
			}

		}
	}
}

// Processes text messages and commands from users.
func handleUpdateMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch {
	case update.Message.Text == types.CommandMetricUnits:
		s.SetSystem(update.Message.Chat.ID, true)
		telegram.SendMessage(bot, update.Message.Chat.ID, types.MetrikUnitOn)
	case update.Message.Text == types.CommandNonMetricUnits:
		s.SetSystem(update.Message.Chat.ID, false)
		telegram.SendMessage(bot, update.Message.Chat.ID, types.MetrikUnitOff)
	case update.Message.Text == types.CommandStart:
		n := update.SentFrom()
		greet := fmt.Sprintf("%s%s%s%s", types.WelcomeMessage, n.FirstName, types.WelcomeMessageEnd, types.HelpMessage)
		telegram.SendMessage(bot, update.Message.Chat.ID, greet)
	case update.Message.Text == types.CommandHelp:
		telegram.SendMessage(bot, update.Message.Chat.ID, types.HelpMessage)
	default:
		s.SetCity(update.Message.Chat.ID, update.Message.Text)
		err = telegram.SendMessageWithInlineKeyboard(bot, update.Message.Chat.ID, types.ChooseOptionMessage, types.CommandCurrent, types.CommandForecast)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
	}
}

// Processes location messages from users.
func handleLocationMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	uLat, uLon := fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
	s.SetLocation(update.Message.Chat.ID, uLat, uLon)
	err = telegram.SendLocationOptions(bot, update.Message.Chat.ID, uLat, uLon)
	if err != nil {
		logger.ForErrorPrint(e.Wrap("", err))
	}
}

// Processes callback queries from users.
func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch {
	case update.CallbackQuery.Data == types.CommandCurrent || update.CallbackQuery.Data == types.CommandForecast:
		if s.GetCity(update.CallbackQuery.Message.Chat.ID) == "" { //s.Data[update.CallbackQuery.Message.Chat.ID].City == "" {
			userMessage = types.MissingCityMessage
			telegram.SendMessage(bot, update.CallbackQuery.Message.Chat.ID, userMessage)
		} else {
			weatherUrl, err := weather.WeatherUrlByCity(s.GetCity(update.CallbackQuery.Message.Chat.ID), tWeather, update.CallbackQuery.Data, s.GetSystem(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			}
			s.SetLast(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			userMessage, err = weather.GetWeather(weatherUrl, update.CallbackQuery.Data, s.GetSystem(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
				userMessage = e.Wrap("", err).Error()
				telegram.SendMessage(bot, update.CallbackQuery.Message.Chat.ID, userMessage)
			} else {
				err = telegram.SendMessageWithInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
				}
			}
		}
	case update.CallbackQuery.Data == types.CommandForecastLocation || update.CallbackQuery.Data == types.CommandCurrentLocation:
		if s.GetLat(update.CallbackQuery.Message.Chat.ID) == "" && s.GetLon(update.CallbackQuery.Message.Chat.ID) == "" {
			userMessage = types.NoLocationProvidedMessage
			telegram.SendMessage(bot, update.CallbackQuery.Message.Chat.ID, userMessage)
		} else {
			weatherUrl, err := weather.WeatherUrlByLocation(s.GetLat(update.CallbackQuery.Message.Chat.ID), s.GetLon(update.CallbackQuery.Message.Chat.ID), tWeather, update.CallbackQuery.Data, s.GetSystem(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			}
			s.SetLast(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			userMessage, err = weather.GetWeather(weatherUrl, update.CallbackQuery.Data, s.GetSystem(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			} else {
				err = telegram.SendMessageWithInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
				}
			}
		}
	}
}

// Processes the "repeat last" callback query, sends the last weather data.
func handleLast(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch {
	// If the users last requested weather type is empty due to a bot restart.
	case !s.Exists(update.CallbackQuery.Message.Chat.ID):
		name := update.SentFrom()
		telegram.SendMessage(bot, update.CallbackQuery.Message.Chat.ID, types.LastDataUnavailable+name.FirstName+types.LastDataUnavailableEnd)
	case s.GetLast(update.CallbackQuery.Message.Chat.ID) == types.CommandCurrent || s.GetLast(update.CallbackQuery.Message.Chat.ID) == types.CommandForecast:
		weatherUrl, err := weather.WeatherUrlByCity(s.GetCity(update.CallbackQuery.Message.Chat.ID), tWeather, s.GetLast(update.CallbackQuery.Message.Chat.ID), s.GetSystem(update.CallbackQuery.Message.Chat.ID))
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
		userMessage, err = weather.GetWeather(weatherUrl, s.GetLast(update.CallbackQuery.Message.Chat.ID), s.GetSystem(update.CallbackQuery.Message.Chat.ID))
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
			userMessage = e.Wrap("", err).Error()
		}
		err = telegram.SendMessageWithInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
	case s.GetLast(update.CallbackQuery.Message.Chat.ID) == types.CommandForecastLocation || s.GetLast(update.CallbackQuery.Message.Chat.ID) == types.CommandCurrentLocation:
		weatherUrl, err := weather.WeatherUrlByLocation(s.GetLat(update.CallbackQuery.Message.Chat.ID), s.GetLon(update.CallbackQuery.Message.Chat.ID), tWeather, s.GetLast(update.CallbackQuery.Message.Chat.ID), s.GetSystem(update.CallbackQuery.Message.Chat.ID))
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
		userMessage, err = weather.GetWeather(weatherUrl, s.GetLast(update.CallbackQuery.Message.Chat.ID), s.GetSystem(update.CallbackQuery.Message.Chat.ID))
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
		err = telegram.SendMessageWithInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
	}
}
