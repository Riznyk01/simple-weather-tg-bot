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

		if update.Message.Text == "/start" {
			userMessage = "Hello! This bot will send you weather information from openweathermap.org. " +
				"Type the name of the city in any language. Use /weather for current weather and /forecast for a 5-day forecast."
		} else if strings.HasPrefix(update.Message.Text, "/weather") {
			city := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/weather"))
			userMessage, err = weather.GetWeather(city, tWeather)
		} else if strings.HasPrefix(update.Message.Text, "/forecast") {
			city := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/forecast"))
			userMessage, err = weather.Get5DayForecast(city, tWeather)
		} else {
			userMessage = "Invalid command. Use /weather [city] for current weather or /forecast [city] for a 5-day forecast."
		}
		fmt.Println(userMessage)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, userMessage)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err = bot.Send(msg)
		if err != nil {
			errorMessage := err.Error()
			log.Println("Error: ", errorMessage)
		}
	}
}
