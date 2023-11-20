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

var city, latStr, lonStr, weatherUrl, userMessage string

func main() {

	err := godotenv.Load(".env.dev")
	if err != nil {
		logger.ForError(err)
	}
	t := os.Getenv("BOT_TOKEN")
	tWeather := os.Getenv("WEATHER_KEY")

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

		if update.Message != nil {
			switch {
			case update.Message.Text == types.CommandStart:
				sendMessage(bot, update.Message.Chat.ID, types.WelcomeMessage+types.HelpMessage)
			case update.Message.Text == types.CommandHelp:
				sendMessage(bot, update.Message.Chat.ID, types.HelpMessage)
			case update.Message.Location != nil:
				latStr, lonStr = fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
				err = sendLocationOptions(bot, update.Message.Chat.ID, latStr, lonStr)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
				}
			default:
				city = update.Message.Text
				err = sendMessageWithInlineKeyboard(bot, update.Message.Chat.ID, types.ChooseOptionMessage, types.CommandCurrent, types.CommandForecast)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
				}
			}
		}

		if update.Message == nil && update.CallbackQuery != nil {
			switch {
			case update.CallbackQuery.Data == types.CommandCurrent || update.CallbackQuery.Data == types.CommandForecast:
				if city != "" {
					weatherUrl, err = weather.WeatherUrlByCity(city, tWeather, update.CallbackQuery.Data)
					if err != nil {
						logger.ForErrorPrint(e.Wrap("", err))
					}
					userMessage, err = weather.GetWeather(weatherUrl, update.CallbackQuery.Data)
					if err != nil {
						logger.ForErrorPrint(e.Wrap("", err))
						userMessage = e.Wrap("", err).Error()
					}
					sendMessage(bot, update.CallbackQuery.Message.Chat.ID, userMessage)
				} else {
					sendMessage(bot, update.CallbackQuery.Message.Chat.ID, types.MissingCityMessage)
				}
			case update.CallbackQuery.Data == types.CommandForecastLocation || update.CallbackQuery.Data == types.CommandCurrentLocation:

				if latStr != "" && lonStr != "" {
					weatherUrl, err = weather.WeatherUrlByLocation(latStr, lonStr, tWeather, update.CallbackQuery.Data)
					if err != nil {
						logger.ForErrorPrint(e.Wrap("", err))
					}
					userMessage, err = weather.GetWeather(weatherUrl, update.CallbackQuery.Data)
					if err != nil {
						logger.ForErrorPrint(e.Wrap("", err))
					}
				} else {
					userMessage = types.NoLocationProvidedMessage
				}
				sendMessage(bot, update.CallbackQuery.Message.Chat.ID, userMessage)
			}
		}
	}
}
