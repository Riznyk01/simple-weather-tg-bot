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
			log.Println("User message:", update.Message.Text, " User's location:", update.Message.Location)
			switch {
			case update.Message.Text != "/start" && update.Message.Text != "/help" && update.Message.Location == nil && update.Message.Text != "current" && update.Message.Text != "5-days forecast" && update.Message.Text != "5-days forecast 📍" && update.Message.Text != "current 📍":
				log.Println("Text received:", update.Message.Text)
				city = update.Message.Text
				userMessage = "Choose an action:"
				err := sendMessageWithKeyboard(bot, update.Message.Chat.ID, userMessage, "current", "5-days forecast")
				if err != nil {
					errorMessage := err.Error()
					log.Println("Error: ", errorMessage)
				}
			case update.Message.Text == "/start":
				userMessage = "Hello! This bot will send you weather information from openweathermap.org. " +
					"Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
			case update.Message.Text == "/help":
				userMessage = "Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
			case update.Message.Text == "current":
				if city != "" {
					weatherUrl := weather.WeatherUrlByCity(city, tWeather, "current")
					log.Println("Case current (by city) choosed, url:", weatherUrl)
					userMessage, err = weather.GetWeather(weatherUrl, "current")
					if err != nil {
						errorMessage := err.Error()
						log.Println("Error: ", errorMessage)
						userMessage = errorMessage
					}
					city = ""
				} else {
					userMessage = "You didn't enter a city.\nPlease enter a city or send your location,\nand then choose the type of weather."
				}
			case update.Message.Text == "5-days forecast":
				if city != "" {
					weatherUrl := weather.WeatherUrlByCity(city, tWeather, "5d3h")
					log.Println("Case forecast (by city) choosed, url:", weatherUrl)
					userMessage, err = weather.GetWeather(weatherUrl, "5d3h")
					if err != nil {
						errorMessage := err.Error()
						log.Println("Error: ", errorMessage)
						userMessage = errorMessage
					}
					city = ""
				} else {
					userMessage = "You did not enter a city.\nPlease enter a city or send your location,\nand then choose the type of weather."
				}
			case update.Message.Location != nil:
				fmt.Println("Case location")
				latStr, lonStr = fmt.Sprintf("%f", update.Message.Location.Latitude), fmt.Sprintf("%f", update.Message.Location.Longitude)
				err := sendLocationOptions(bot, update.Message.Chat.ID, latStr, lonStr)
				if err != nil {
					errorMessage := err.Error()
					log.Println("Error: ", errorMessage)
				}
			case update.Message.Text == "5-days forecast 📍":
				weatherUrl := weather.WeatherUrlByLocation(latStr, lonStr, tWeather, "5d3h")
				log.Println("5-days forecast (by location) choosed, url:", weatherUrl)
				userMessage, err = weather.GetWeather(weatherUrl, "5d3h")
				if err != nil {
					errorMessage := err.Error()
					log.Println("5-days forecast (by location) error: ", errorMessage)
				}
			case update.Message.Text == "current 📍":
				weatherUrl := weather.WeatherUrlByLocation(latStr, lonStr, tWeather, "current")
				log.Println("Current weather (by location) choosed, url:", weatherUrl)
				userMessage, err = weather.GetWeather(weatherUrl, "current")
				if err != nil {
					errorMessage := err.Error()
					log.Println("Current weather (by location) error: ", errorMessage)
				}
			default:
				userMessage = "Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
			}
			if userMessage != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, userMessage)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID

				_, err = bot.Send(msg)
				if err != nil {
					errorMessage := err.Error()
					log.Println("Error: ", errorMessage)
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
	chooseWeatherType := fmt.Sprintf("Your location: %s %v, %v. Choose an action:", latStr, lonStr, "Choose an action:")
	return sendMessageWithKeyboard(bot, chatID, chooseWeatherType, "5-days forecast 📍", "current 📍")
}
