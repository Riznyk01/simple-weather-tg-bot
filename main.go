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

var city, latStr, lonStr, weatherUrl, userMessage, tWeather string
var err error
var userUnits = make(map[int64]bool)
var units bool

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
			go handleCallbackQuery(bot, update)
		}
	}
}
func handleUpdateMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch {
	case update.Message.Text == types.CommandMetricUnits:
		setUserUnitsToMetric(update.Message.Chat.ID, true)
		sendMessage(bot, update.Message.Chat.ID, types.MetrikUnitOn)
	case update.Message.Text == types.CommandNonMetricUnits:
		setUserUnitsToMetric(update.Message.Chat.ID, false)
		sendMessage(bot, update.Message.Chat.ID, types.MetrikUnitOff)
	case update.Message.Text == types.CommandStart:
		sendMessage(bot, update.Message.Chat.ID, types.WelcomeMessage+types.HelpMessage)
	case update.Message.Text == types.CommandHelp:
		sendMessage(bot, update.Message.Chat.ID, types.HelpMessage)
	default:
		city = update.Message.Text
		err = sendMessageWithInlineKeyboard(bot, update.Message.Chat.ID, types.ChooseOptionMessage, types.CommandCurrent, types.CommandForecast)
		if err != nil {
			logger.ForErrorPrint(e.Wrap("", err))
		}
	}
}
func handleLocationMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	latStr, lonStr = fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
	err = sendLocationOptions(bot, update.Message.Chat.ID, latStr, lonStr)
	if err != nil {
		logger.ForErrorPrint(e.Wrap("", err))
	}
}
func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch {
	case update.CallbackQuery.Data == types.CommandCurrent || update.CallbackQuery.Data == types.CommandForecast:
		if city == "" {
			userMessage = types.MissingCityMessage
		} else {
			weatherUrl, units, err = weather.WeatherUrlByCity(city, tWeather, update.CallbackQuery.Data, getUserUnitsToMetric(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			}
			userMessage, err = weather.GetWeather(weatherUrl, update.CallbackQuery.Data, units)
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
				userMessage = e.Wrap("", err).Error()
			}
		}
	case update.CallbackQuery.Data == types.CommandForecastLocation || update.CallbackQuery.Data == types.CommandCurrentLocation:
		if latStr == "" && lonStr == "" {
			userMessage = types.NoLocationProvidedMessage
		} else {
			weatherUrl, units, err = weather.WeatherUrlByLocation(latStr, lonStr, tWeather, update.CallbackQuery.Data, getUserUnitsToMetric(update.CallbackQuery.Message.Chat.ID))
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			}
			userMessage, err = weather.GetWeather(weatherUrl, update.CallbackQuery.Data, units)
			if err != nil {
				logger.ForErrorPrint(e.Wrap("", err))
			}
		}
	}
	sendMessage(bot, update.CallbackQuery.Message.Chat.ID, userMessage)
}

func setUserUnitsToMetric(chatID int64, metric bool) {
	userUnits[chatID] = metric
}
func getUserUnitsToMetric(chatID int64) bool {
	return userUnits[chatID]
}
