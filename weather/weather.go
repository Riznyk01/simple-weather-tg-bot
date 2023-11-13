package weather

import (
	"SimpleWeatherTgBot/types"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetWeather(city, tWeather string) (types.WeatherResponse, error) {

	weatherUrl := "https://api.openweathermap.org/data/2.5/weather?"

	u, err := url.Parse(weatherUrl)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		return types.WeatherResponse{}, fmt.Errorf("error: %s", errorMessage)
	}
	q := url.Values{}
	q.Add("q", city)
	q.Add("appid", tWeather)
	q.Add("units", "metric")
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	//fmt.Println(fullUrlGet)
	resp, err := http.Get(fullUrlGet)
	if err != nil {
		return types.WeatherResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		return types.WeatherResponse{}, fmt.Errorf("error: %s", errorMessage)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			var errorResponse struct {
				Cod     string `json:"cod"`
				Message string `json:"message"`
			}
			err = json.Unmarshal(body, &errorResponse)
			if err == nil {
				return types.WeatherResponse{}, fmt.Errorf("%s"+". Try another city name.", errorResponse.Message)
			}
		}
		return types.WeatherResponse{}, fmt.Errorf("Failed to get weather data. Status code: %d", resp.StatusCode)
	}

	var weatherResponse types.WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		return types.WeatherResponse{}, fmt.Errorf("error: %s", errorMessage)
	}
	return weatherResponse, nil
}
