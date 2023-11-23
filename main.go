package main

import (
	"SimpleWeatherTgBot/lib/e"
	"SimpleWeatherTgBot/logger"
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
	userData              = make(map[int64]types.UserData)
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
		currentData := userData[update.Message.Chat.ID]
		currentData.Metric = true
		userData[update.Message.Chat.ID] = currentData
		sendMessage(bot, update.Message.Chat.ID, types.MetrikUnitOn)
	case update.Message.Text == types.CommandNonMetricUnits:
		currentData := userData[update.Message.Chat.ID]
		currentData.Metric = false
		userData[update.Message.Chat.ID] = currentData
		sendMessage(bot, update.Message.Chat.ID, types.MetrikUnitOff)
	case update.Message.Text == types.CommandStart:
		sendMessage(bot, update.Message.Chat.ID, types.WelcomeMessage+types.HelpMessage)
	case update.Message.Text == types.CommandHelp:
		sendMessage(bot, update.Message.Chat.ID, types.HelpMessage)
	default:
		currentData := userData[update.Message.Chat.ID]
		currentData.City = update.Message.Text
		userData[update.Message.Chat.ID] = currentData
		err = sendMessageWithInlineKeyboard(bot, update.Message.Chat.ID, types.ChooseOptionMessage, types.CommandCurrent, types.CommandForecast)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
	}
}

// Processes location messages from users.
func handleLocationMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	currentData := userData[update.Message.Chat.ID]
	currentData.Lat = fmt.Sprintf("%f", update.Message.Location.Latitude)
	currentData.Lon = fmt.Sprintf("%f", update.Message.Location.Longitude)
	userData[update.Message.Chat.ID] = currentData
	err = sendLocationOptions(bot, update.Message.Chat.ID, userData[update.Message.Chat.ID].Lat, userData[update.Message.Chat.ID].Lon)
	if err != nil {
		logger.ForErrorPrint(e.Wrap("", err))
	}
}

// Processes callback queries from users.
func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch {
	case update.CallbackQuery.Data == types.CommandCurrent || update.CallbackQuery.Data == types.CommandForecast:
		if userData[update.CallbackQuery.Message.Chat.ID].City == "" {
			userMessage = types.MissingCityMessage
		} else {
			currentData := userData[update.CallbackQuery.Message.Chat.ID]
			weatherUrl, err := weather.WeatherUrlByCity(userData[update.CallbackQuery.Message.Chat.ID].City, tWeather, update.CallbackQuery.Data, userData[update.CallbackQuery.Message.Chat.ID].Metric)
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			}
			currentData.Last = update.CallbackQuery.Data
			userData[update.CallbackQuery.Message.Chat.ID] = currentData

			userMessage, err = weather.GetWeather(weatherUrl, update.CallbackQuery.Data, userData[update.CallbackQuery.Message.Chat.ID].Metric)
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
				userMessage = e.Wrap("", err).Error()
				sendMessage(bot, update.CallbackQuery.Message.Chat.ID, userMessage)
			} else {
				err = sendMessageWithInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
				}
			}
		}
	case update.CallbackQuery.Data == types.CommandForecastLocation || update.CallbackQuery.Data == types.CommandCurrentLocation:
		if userData[update.CallbackQuery.Message.Chat.ID].Lat == "" && userData[update.CallbackQuery.Message.Chat.ID].Lon == "" {
			userMessage = types.NoLocationProvidedMessage
		} else {
			currentData := userData[update.CallbackQuery.Message.Chat.ID]

			weatherUrl, err := weather.WeatherUrlByLocation(userData[update.CallbackQuery.Message.Chat.ID].Lat, userData[update.CallbackQuery.Message.Chat.ID].Lon, tWeather, update.CallbackQuery.Data, userData[update.CallbackQuery.Message.Chat.ID].Metric)
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			}

			currentData.Last = update.CallbackQuery.Data
			userData[update.CallbackQuery.Message.Chat.ID] = currentData
			userMessage, err = weather.GetWeather(weatherUrl, update.CallbackQuery.Data, userData[update.CallbackQuery.Message.Chat.ID].Metric)
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			}
		}
		err = sendMessageWithInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
	}
}

// Processes the "repeat last" callback query, sends the last weather data.
func handleLast(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch {
	case userData[update.CallbackQuery.Message.Chat.ID].Last == "":
		sendMessage(bot, update.CallbackQuery.Message.Chat.ID, types.LastDataUnavailable)
	case userData[update.CallbackQuery.Message.Chat.ID].Last == types.CommandCurrent || userData[update.CallbackQuery.Message.Chat.ID].Last == types.CommandForecast:
		weatherUrl, err := weather.WeatherUrlByCity(userData[update.CallbackQuery.Message.Chat.ID].City, tWeather, userData[update.CallbackQuery.Message.Chat.ID].Last, userData[update.CallbackQuery.Message.Chat.ID].Metric)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
		userMessage, err = weather.GetWeather(weatherUrl, userData[update.CallbackQuery.Message.Chat.ID].Last, userData[update.CallbackQuery.Message.Chat.ID].Metric)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
			userMessage = e.Wrap("", err).Error()
		}
		err = sendMessageWithInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
	case userData[update.CallbackQuery.Message.Chat.ID].Last == types.CommandForecastLocation || userData[update.CallbackQuery.Message.Chat.ID].Last == types.CommandCurrentLocation:
		weatherUrl, err := weather.WeatherUrlByLocation(userData[update.CallbackQuery.Message.Chat.ID].Lat, userData[update.CallbackQuery.Message.Chat.ID].Lon, tWeather, userData[update.CallbackQuery.Message.Chat.ID].Last, userData[update.CallbackQuery.Message.Chat.ID].Metric)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
		userMessage, err = weather.GetWeather(weatherUrl, userData[update.CallbackQuery.Message.Chat.ID].Last, userData[update.CallbackQuery.Message.Chat.ID].Metric)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
		err = sendMessageWithInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, userMessage, types.CommandLast)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
	}
}
