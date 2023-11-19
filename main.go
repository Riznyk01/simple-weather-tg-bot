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

var city, latStr, lonStr string

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
			case update.Message.Text == types.CommandCurrent || update.Message.Text == types.CommandForecast:
				if city != "" {
					weatherUrl, err := weather.WeatherUrlByCity(city, tWeather, update.Message.Text)
					if err != nil {
						logger.ForErrorPrint(e.Wrap("", err))
					}
					userMessage, err := weather.GetWeather(weatherUrl, update.Message.Text)
					if err != nil {
						logger.ForErrorPrint(e.Wrap("", err))
						userMessage = e.Wrap("", err).Error()
					}
					sendMessage(bot, update.Message.Chat.ID, userMessage)
					city = ""
				} else {
					sendMessage(bot, update.Message.Chat.ID, types.MissingCityMessage)
				}
			case update.Message.Location != nil:
				latStr, lonStr = fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
				err = sendLocationOptions(bot, update.Message.Chat.ID, latStr, lonStr)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
				}
			case update.Message.Text == types.CommandForecastLocation || update.Message.Text == types.CommandCurrentLocation:
				weatherUrl, err := weather.WeatherUrlByLocation(latStr, lonStr, tWeather, update.Message.Text)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
				}
				userMessage, err := weather.GetWeather(weatherUrl, update.Message.Text)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
					userMessage = e.Wrap("", err).Error()
				}
				sendMessage(bot, update.Message.Chat.ID, userMessage)
			default:
				city = update.Message.Text
				err = sendMessageWithKeyboard(bot, update.Message.Chat.ID, types.ChooseOptionMessage, types.CommandCurrent, types.CommandForecast)
				if err != nil {
					logger.ForErrorPrint(e.Wrap("", err))
				}
			}
		}
	}
}
