package main

import (
	"SimpleWeatherTgBot/utils"
	"SimpleWeatherTgBot/weather"
	"fmt"
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

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {

			if update.Message.Text == "/start" {
				userMessage = "Hello! This bot will send you weather information from openweathermap.org in response to your message with the name of the city in any language. \nSimply enter the city name and send it to the bot."
			} else {

				weatherData, err := weather.GetWeather(update.Message.Text, tWeather)
				if err != nil {
					return
				}
				if weatherData.Weather[0].Main == "Rain" {
					weatherData.Weather[0].Main = "🌧 Rain"
				} else if weatherData.Weather[0].Main == "Clouds" {
					weatherData.Weather[0].Main = "☁️ Clouds"
				} else if weatherData.Weather[0].Main == "Clear" {
					weatherData.Weather[0].Main = "✨ Clear"
				}
				userMessage = fmt.Sprintf("%s %s - %s 🌡 %.2f°C 💧 %d%%\n\nFeelsLike %.2f°C  🔺 %.2f°C ️ 🔻 %.2f°C \n%d hPa / %.2f mmHg\n %.2f m/s / %s \n\n🌅  %s 🌉  %s",
					weatherData.Sys.Country,
					weatherData.Name,
					weatherData.Weather[0].Main,
					weatherData.Main.Temp,
					weatherData.Main.Humidity,
					weatherData.Main.FeelsLike,
					weatherData.Main.TempMax,
					weatherData.Main.TempMin,
					weatherData.Main.Pressure,
					utils.HPaToMmHg(float64(weatherData.Main.Pressure)),
					weatherData.Wind.Speed,
					utils.DegreesToDirection(weatherData.Wind.Deg),
					utils.TimeStampToHuman(weatherData.Sys.Sunrise, weatherData.Timezone),
					utils.TimeStampToHuman(weatherData.Sys.Sunset, weatherData.Timezone))
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, userMessage)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(msg)
			if err != nil {
				return
			}
		}
	}
}
