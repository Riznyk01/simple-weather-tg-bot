package weather

import (
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Response(fullUrl, tWeather string, update types.Update) error {
	var respMessage types.RespMessage
	respMessage.ChatId = update.Message.Chat.ChatId

	if update.Message.Text == "/start" {
		respMessage.Text = "Hello! This bot will send you weather information from openweathermap.org in response to your message with the name of the city in any language. \nSimply enter the city name and send it to the bot."
	} else {

		weatherData, err := GetWeather(update.Message.Text, tWeather)
		if err != nil {
			return err
		}
		if weatherData.Weather[0].Main == "Rain" {
			weatherData.Weather[0].Main = "🌧 Rain"
		} else if weatherData.Weather[0].Main == "Clouds" {
			weatherData.Weather[0].Main = "☁️ Clouds"
		} else if weatherData.Weather[0].Main == "Clear" {
			weatherData.Weather[0].Main = "✨ Clear"
		}
		respMessage.Text = fmt.Sprintf("%s %s - %s 🌡 %.2f°C 💧 %d%%\n\nFeelsLike %.2f°C  🔺 %.2f°C ️ 🔻 %.2f°C \n%d hPa / %.2f mmHg\n %.2f m/s / %s \n\n🌅  %s 🌉  %s",
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
	//}
	buf, err := json.Marshal(respMessage)
	if err != nil {
		return err
		//log.Println("Smth went wrong: ", err.Error())
	}
	_, err = http.Post(fullUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
		//log.Println("Smth went wrong: ", err.Error())
	}
	return nil
}
