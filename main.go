package main

import (
	"SimpleWeatherTgBot/weather"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var city, latStr, lonStr, userMessage string

// Constants for commands
const (
	CommandStart            = "/start"
	CommandHelp             = "/help"
	CommandCurrent          = "current"
	CommandForecast         = "5-days forecast"
	CommandForecastLocation = "5-days forecast 📍"
	CommandCurrentLocation  = "current 📍"
)

// Constants for messages
const (
	WelcomeMessage      = "Hello! This bot will send you weather information from openweathermap.org. "
	HelpMessage         = "Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
	MissingCityMessage  = "You didn't enter a city.\nPlease enter a city or send your location,\nand then choose the type of weather."
	ChooseOptionMessage = "Choose an action:"
)

func main() {

	err := godotenv.Load(".env.dev")
	if err != nil {
		return
	}
	t := os.Getenv("BOT_TOKEN")
	tWeather := os.Getenv("WEATHER_KEY")

	bot, err := tgbotapi.NewBotAPI(t)
	if err != nil {
		log.Panic(err)
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
			case update.Message.Text == CommandStart:
				sendMessage(bot, update.Message.Chat.ID, WelcomeMessage+HelpMessage)
			case update.Message.Text == CommandHelp:
				sendMessage(bot, update.Message.Chat.ID, HelpMessage)
			case update.Message.Text == CommandCurrent || update.Message.Text == CommandForecast:
				if city != "" {
					weatherUrl := weather.WeatherUrlByCity(city, tWeather, update.Message.Text)
					userMessage, err = weather.GetWeather(weatherUrl, update.Message.Text)
					if err != nil {
						userMessage = HandleErrorMessage("", err)
					}
					sendMessage(bot, update.Message.Chat.ID, userMessage)
					city = ""
				} else {
					sendMessage(bot, update.Message.Chat.ID, MissingCityMessage)
				}
			case update.Message.Location != nil:
				latStr, lonStr = fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
				err = sendLocationOptions(bot, update.Message.Chat.ID, latStr, lonStr)
				if err != nil {
					HandleError("", err)
				}
			case update.Message.Text == CommandForecastLocation || update.Message.Text == CommandCurrentLocation:
				weatherUrl := weather.WeatherUrlByLocation(latStr, lonStr, tWeather, update.Message.Text)
				userMessage, err = weather.GetWeather(weatherUrl, update.Message.Text)
				if err != nil {
					HandleError("", err)
				}
				sendMessage(bot, update.Message.Chat.ID, userMessage)
			default:
				city = update.Message.Text
				err = sendMessageWithKeyboard(bot, update.Message.Chat.ID, ChooseOptionMessage, CommandCurrent, CommandForecast)
				if err != nil {
					HandleError("", err)
				}
			}
		}
	}
}
