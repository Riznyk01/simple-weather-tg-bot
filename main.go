package main

import (
	"SimpleWeatherTgBot/weather"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	var userMessage string

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
			if update.Message.Text == "/start" {
				userMessage = "Hello! This bot will send you weather information from openweathermap.org in response to your message with the name of the city in any language. \nSimply enter the city name and send it to the bot."
			} else {
				userMessage, err = weather.GetWeather(update.Message.Text, tWeather)
				if err != nil {
					errorMessage := err.Error()
					log.Println("Error getting weather data: ", errorMessage)
					userMessage = "Error getting weather data: " + errorMessage
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, userMessage)
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(msg)
			if err != nil {
				errorMessage := err.Error()
				log.Println("Error: ", errorMessage)
			}
		}
	}
}
