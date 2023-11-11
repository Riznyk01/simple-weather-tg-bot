package main

import (
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

	offset := 0
	for {
		updates, err := getUpdates(fullUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = response(directGeoUrl, endOfDirectGeoUrl, fullUrl, update)
			offset = update.UpdateId + 1
		}
		//fmt.Println(updates)
	}
}

func getUpdates(fullUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(fullUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func response(weatherApi, endOfWeatherUrl, fullUrl string, update Update) error {
	var respMessage RespMessage
	respMessage.ChatId = update.Message.Chat.ChatId

	if update.Message.Text == "/start" {
		respMessage.Text = "Hello, this bot will send you weather from openweathermap.org in response to your message with the name of the city in any language."
	} else {
		//city:=update.Message.Text
		geo, err := CoordinatesByLocationName(weatherApi, endOfWeatherUrl, update)
		if err != nil {
			return err
		}
		if len(geo) != 0 {
			latStr := strconv.FormatFloat(geo[0].Lat, 'f', -1, 64)
			lonStr := strconv.FormatFloat(geo[0].Lon, 'f', -1, 64)
			fmt.Println(latStr, lonStr)
		}
		//respMessage.Text =

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

func getWeather(fullUrl string, offset int) ([]Update, error) {
	//https://api.openweathermap.org/data/3.0/onecall?lat={lat}&lon={lon}&exclude={part}&appid={API key}
	resp, err := http.Get(fullUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func CoordinatesByLocationName(weatherApi string, endOfWeatherUrl string, update Update) ([]Geocoding, error) {
	city := update.Message.Text
	resp, err := http.Get(weatherApi + city + endOfWeatherUrl)

	if err != nil {
		fmt.Println(err)
		return []Geocoding{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return []Geocoding{}, err
	}
	var geocoding []Geocoding
	err = json.Unmarshal(body, &geocoding)
	if err != nil {
		fmt.Println(err)
		return []Geocoding{}, err
	}
	return geocoding, nil
}
