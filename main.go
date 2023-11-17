package main

import (
	"SimpleWeatherTgBot/weather"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var city, latStr, lonStr string

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

// Constants for weather types
const (
	WeatherTypeCurrent  = "current"
	WeatherTypeForecast = "5d3h"
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
			var userMessage string
			var err error
			switch {
			case update.Message.Text == CommandStart:
				sendMessage(bot, update.Message.Chat.ID, WelcomeMessage+HelpMessage)
			case update.Message.Text == CommandHelp:
				sendMessage(bot, update.Message.Chat.ID, HelpMessage)
			case update.Message.Text == CommandCurrent:
				if city != "" {
					weatherUrl := weather.WeatherUrlByCity(city, tWeather, WeatherTypeCurrent)
					userMessage, err = weather.GetWeather(weatherUrl, WeatherTypeCurrent)
					if err != nil {
						userMessage = HandleErrorMessage("", err)
					}
					sendMessage(bot, update.Message.Chat.ID, userMessage)
					city = ""
				} else {
					sendMessage(bot, update.Message.Chat.ID, MissingCityMessage)
				}
			case update.Message.Text == CommandForecast:
				if city != "" {
					weatherUrl := weather.WeatherUrlByCity(city, tWeather, WeatherTypeForecast)
					userMessage, err = weather.GetWeather(weatherUrl, WeatherTypeForecast)
					if err != nil {
						userMessage = HandleErrorMessage("", err)
					}
					sendMessage(bot, update.Message.Chat.ID, userMessage)
					city = ""
				} else {
					sendMessage(bot, update.Message.Chat.ID, MissingCityMessage)
				}
			case update.Message.Location != nil:
				latStr = fmt.Sprintf("%f", update.Message.Location.Latitude)
				lonStr = fmt.Sprintf("%f", update.Message.Location.Longitude)
				err = sendLocationOptions(bot, update.Message.Chat.ID, latStr, lonStr)
				if err != nil {
					HandleError("", err)
				}
			case update.Message.Text == CommandForecastLocation:
				weatherUrl := weather.WeatherUrlByLocation(latStr, lonStr, tWeather, WeatherTypeForecast)
				userMessage, err = weather.GetWeather(weatherUrl, WeatherTypeForecast)
				if err != nil {
					HandleError("", err)
				}
				sendMessage(bot, update.Message.Chat.ID, userMessage)
			case update.Message.Text == CommandCurrentLocation:
				weatherUrl := weather.WeatherUrlByLocation(latStr, lonStr, tWeather, WeatherTypeCurrent)
				userMessage, err = weather.GetWeather(weatherUrl, WeatherTypeCurrent)
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

// sendMessageWithKeyboard sends a message with the specified text and keyboard buttons.
func sendMessageWithKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string, buttons ...string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	keyboardButtons := make([]tgbotapi.KeyboardButton, len(buttons))
	for i, button := range buttons {
		keyboardButtons[i] = tgbotapi.NewKeyboardButton(button)
	}
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(keyboardButtons...))
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	return err
}

// sendLocationOptions sends a message with location-related options.
func sendLocationOptions(bot *tgbotapi.BotAPI, chatID int64, latStr, lonStr string) error {
	chooseWeatherType := fmt.Sprintf("Your location: %s, %v\n%s", latStr, lonStr, ChooseOptionMessage)
	return sendMessageWithKeyboard(bot, chatID, chooseWeatherType, CommandForecastLocation, CommandCurrentLocation)
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err := bot.Send(msg)
	if err != nil {
		HandleError("", err)
	}
}
