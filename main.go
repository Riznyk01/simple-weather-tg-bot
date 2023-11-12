package main

import (
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/utils"
	"SimpleWeatherTgBot/weather"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		return
	}
	t := os.Getenv("BOT_TOKEN")

	botApi := "https://api.telegram.org/bot"
	fullUrl := botApi + t

	tWeather := os.Getenv("WEATHER_KEY")

	directGeoUrl := "http://api.openweathermap.org/geo/1.0/direct?q="
	limit := "1"
	endOfDirectGeoUrl := "&limit=" + limit + "&appid=" + tWeather

	weatherUrl := "https://api.openweathermap.org/data/2.5/weather?"

	offset := 0
	for {
		updates, err := getUpdates(fullUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = response(weatherUrl, directGeoUrl, endOfDirectGeoUrl, fullUrl, tWeather, update)
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(fullUrl string, offset int) ([]types.Update, error) {
	resp, err := http.Get(fullUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse types.RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func response(weatherUrl, directGeoUrl, endOfDirectGeoUrl, fullUrl, tWeather string, update types.Update) error {
	var respMessage types.RespMessage
	respMessage.ChatId = update.Message.Chat.ChatId

	if update.Message.Text == "/start" {
		respMessage.Text = "Hello, this bot will send you weather from openweathermap.org in response to your message with the name of the city in any language."
	} else {
		geo, err := weather.CoordinatesByLocationName(directGeoUrl, endOfDirectGeoUrl, update)
		if err != nil {
			return err
		}
		if len(geo) != 0 {
			latStr := strconv.FormatFloat(geo[0].Lat, 'f', -1, 64)
			lonStr := strconv.FormatFloat(geo[0].Lon, 'f', -1, 64)

			weatherData, err := getWeather(weatherUrl, latStr, lonStr, tWeather)
			if err != nil {
				return err
			}
			fmt.Println(weatherData)
			if weatherData.Weather[0].Main == "Rain" {
				weatherData.Weather[0].Main = " 🌧"
			} else if weatherData.Weather[0].Main == "Clouds" {
				weatherData.Weather[0].Main += " ☁️"
			}

			respMessage.Text = fmt.Sprintf("%s %s - %s \n\n🌡Now %.2f°C     FeelsLike %.2f°C\n       Max %.2f°C     ️Min %.2f°C 💧 %d%%\n\n 💨%d hPa / %.2f mmHg\n        %.2f m/s / %s \n\n🌅  %s\n🌉  %s",
				weatherData.Sys.Country,
				weatherData.Name,
				weatherData.Weather[0].Main,
				weatherData.Main.Temp,
				weatherData.Main.FeelsLike,
				weatherData.Main.TempMax,
				weatherData.Main.TempMin,
				weatherData.Main.Humidity,
				weatherData.Main.Pressure,
				utils.HPaToMmHg(float64(weatherData.Main.Pressure)),
				weatherData.Wind.Speed,
				utils.DegreesToDirection(weatherData.Wind.Deg),
				utils.TimeStampToHuman(weatherData.Sys.Sunrise, weatherData.Timezone),
				utils.TimeStampToHuman(weatherData.Sys.Sunset, weatherData.Timezone))
		}
	}
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

func getWeather(weatherUrl, latStr, lonStr, tWeather string) (types.WeatherResponse, error) {
	resp, err := http.Get(weatherUrl + "lat=" + latStr + "&lon=" + lonStr + "&appid=" + tWeather + "&units=metric")
	if err != nil {
		return types.WeatherResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("resp error:", err)
		return types.WeatherResponse{}, err
	}
	var weatherResponse types.WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		fmt.Println("getWeather func err:", err)
		return types.WeatherResponse{}, err
	}
	return weatherResponse, nil
}
