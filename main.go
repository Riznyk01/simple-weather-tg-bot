package main

import (
	"SimpleWeatherTgBot/weather"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var lat, lon float64

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

	if update.Message != nil {
		var userMessage string
		var err error

		switch {
		case update.Message.Text == "/f":
			userMessage = "You haven't entered the city name. Please enter it in the following format: /f [city_name]."
		case update.Message.Text == "/w":
			userMessage = "You haven't entered the city name. Please enter it in the following format: /w [city_name]."
		case update.Message.Text == "/start":
			userMessage = "Hello! This bot will send you weather information from openweathermap.org. " +
				"Type the name of the city in any language. Use /w for current weather and /f for a 5-day forecast." +
				"If you have a city with a common name and want weather for a specific location, you can send your location to get accurate weather information."
		case update.Message.Text == "/help":
			userMessage = "/f - command for 5-day weather forecast. Format: /f [city_name]\n /w - command for current weather. Format: /w [city_name]"
		case strings.HasPrefix(update.Message.Text, "/w"):
			city := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/w"))
			weatherNowUrl := weather.WeatherNowUrlByCity(city, tWeather)
			userMessage, err = weather.GetWeather(weatherNowUrl)
		case strings.HasPrefix(update.Message.Text, "/f"):
			city := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/f"))
			weather5d3hUrl := weather.Weather5d3hUrlByCity(city, tWeather)
			userMessage, err = weather.Get5DayForecast(weather5d3hUrl)
		case update.Message.Location != nil:
			lat = update.Message.Location.Latitude
			lon = update.Message.Location.Longitude
			chooseWeatherType := fmt.Sprintf("Your location: %v, %v. Choose an action:", lat, lon)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, chooseWeatherType)
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("5-day forecast"),
					tgbotapi.NewKeyboardButton("current weather"),
				),
			)
			msg.ReplyMarkup = keyboard
			_, err := bot.Send(msg)
			if err != nil {
				errorMessage := err.Error()
				log.Println("Error: ", errorMessage)
			}
		case update.Message.Text == "5-day forecast":
			latStr := fmt.Sprintf("%f", lat)
			lonStr := fmt.Sprintf("%f", lon)
			weatherNowUrl := weather.Weather5d3hUrlByLocation(latStr, lonStr, tWeather)
			userMessage, err = weather.Get5DayForecast(weatherNowUrl)
		case update.Message.Text == "current weather":
			latStr := fmt.Sprintf("%f", lat)
			lonStr := fmt.Sprintf("%f", lon)
			weatherNowUrl := weather.WeatherNowUrlByLocation(latStr, lonStr, tWeather)
			fmt.Println(weatherNowUrl)
			userMessage, err = weather.GetWeather(weatherNowUrl)
		default:
			userMessage = "Invalid command. Use /w [city] for current weather or /f [city] for a 5-day forecast."
		}

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
