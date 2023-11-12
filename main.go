package main

import (
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/weather"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		return
	}
	t := os.Getenv("BOT_TOKEN")
	tWeather := os.Getenv("WEATHER_KEY")

	botApi := "https://api.telegram.org/bot"
	baseUrl := botApi + t

	directGeoUrl := "http://api.openweathermap.org/geo/1.0/direct?q="
	limit := "1"
	endOfDirectGeoUrl := "&limit=" + limit + "&appid=" + tWeather

	weatherUrl := "https://api.openweathermap.org/data/2.5/weather?"

	offset := 0
	for {
		updates, err := getUpdates(baseUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = weather.Response(weatherUrl, directGeoUrl, endOfDirectGeoUrl, baseUrl, tWeather, update)
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(baseUrlGet string, offset int) ([]types.Update, error) {

	u, err := url.Parse(baseUrlGet + "/getUpdates")
	if err != nil {
		fmt.Println("Error parsing URL (getUpdates):", err)
		return nil, err
	}
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()

	resp, err := http.Get(fullUrlGet)
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
