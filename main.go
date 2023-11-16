package main

import (
	"SimpleWeatherTgBot/weather"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var lat, lon float64
var city string

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
		go handleUpdate(bot, update, tWeather)
	}
}
func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, tWeather string) {
	text0 := "Choose an action:"
	text1 := "Hello! This bot will send you weather information from openweathermap.org. "
	text2 := "Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
	text4 := "current"
	text5 := "5-days forecast"
	text6 := "5-days forecast 📍"
	text7 := "current 📍"
	text8 := "Your location:"

	if update.Message != nil {
		var userMessage string
		var err error
		log.Println("User message:", update.Message.Text, " User's location:", update.Message.Location)
		switch {
		case update.Message.Text != "/start" && update.Message.Text != "/help" && update.Message.Location == nil && update.Message.Text != text4 && update.Message.Text != text5 && update.Message.Text != text6 && update.Message.Text != text7:
			log.Println("Text received:", update.Message.Text)
			city = update.Message.Text
			userMessage = text0
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, userMessage)
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(text4),
					tgbotapi.NewKeyboardButton(text5),
				),
			)
			msg.ReplyMarkup = keyboard
			_, err := bot.Send(msg)
			if err != nil {
				errorMessage := err.Error()
				log.Println("Error: ", errorMessage)
			}
		case update.Message.Text == "/start":
			userMessage = text1 + text2
		case update.Message.Text == "/help":
			userMessage = text2
		case update.Message.Text == text4:
			if city != "" {
				weatherUrl := weather.WeatherUrlByCity(city, tWeather, "current")
				log.Println("Case current (by city) choosed, url:", weatherUrl)
				userMessage, err = weather.GetWeather(weatherUrl)
				if err != nil {
					errorMessage := err.Error()
					log.Println("Error: ", errorMessage)
					userMessage = errorMessage
				}
				city = ""
			} else {
				userMessage = "You did not enter a city.\nPlease enter a city or send your location,\nand then choose the type of weather."
			}
		case update.Message.Text == text5:
			if city != "" {
				weatherUrl := weather.WeatherUrlByCity(city, tWeather, "5d3h")
				log.Println("Case forecast (by city) choosed, url:", weatherUrl)
				userMessage, err = weather.Get5DayForecast(weatherUrl)
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
			lat = update.Message.Location.Latitude
			lon = update.Message.Location.Longitude
			chooseWeatherType := fmt.Sprintf("%s %v, %v. %s", text8, lat, lon, text0)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, chooseWeatherType)
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(text6),
					tgbotapi.NewKeyboardButton(text7),
				),
			)
			msg.ReplyMarkup = keyboard
			_, err := bot.Send(msg)
			if err != nil {
				errorMessage := err.Error()
				log.Println("Error: ", errorMessage)
			}
		case update.Message.Text == text6:
			latStr := fmt.Sprintf("%f", lat)
			lonStr := fmt.Sprintf("%f", lon)
			weatherNowUrl := weather.WeatherUrlByLocation(latStr, lonStr, tWeather, "5d3h")
			log.Println("5-days forecast (by location) choosed, url:", weatherNowUrl)
			userMessage, err = weather.Get5DayForecast(weatherNowUrl)
			if err != nil {
				errorMessage := err.Error()
				log.Println("5-days forecast (by location) error: ", errorMessage)
			}
		case update.Message.Text == text7:
			latStr := fmt.Sprintf("%f", lat)
			lonStr := fmt.Sprintf("%f", lon)
			weatherNowUrl := weather.WeatherUrlByLocation(latStr, lonStr, tWeather, "current")
			log.Println("Current weather (by location) choosed, url:", weatherNowUrl)
			userMessage, err = weather.GetWeather(weatherNowUrl)
			if err != nil {
				errorMessage := err.Error()
				log.Println("Current weather (by location) error: ", errorMessage)
			}
		default:
			userMessage = text2
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
