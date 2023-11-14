package weather

import (
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func Get5DayForecast(city, tWeather string) (string, error) {

	weatherUrl := "https://api.openweathermap.org/data/2.5/forecast?"

	u, err := url.Parse(weatherUrl)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		return "", fmt.Errorf("error: %s", errorMessage)
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
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		return "", fmt.Errorf("error: %s", errorMessage)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			var errorResponse struct {
				Cod     string `json:"cod"`
				Message string `json:"message"`
			}
			err = json.Unmarshal(body, &errorResponse)
			if err == nil {
				return "", fmt.Errorf("%s. Try another city name", errorResponse.Message)
			}
		}
		return "", fmt.Errorf("Failed to get weather data. Status code: %d", resp.StatusCode)
	}

	var forecastData types.WeatherResponse5d3h
	err = json.Unmarshal(body, &forecastData)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		return "", fmt.Errorf("error: %s", errorMessage)
	}
	var forecast string

	for _, entry := range forecastData.List {
		hours := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "15")
		monthDay := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "01-02")

		if hours == "01" || hours == "02" {
			forecast += monthDay + "\n"
		}

		forecast += fmt.Sprintf("%s %v %v°C %d%% %.1f mmHg %.1f m/s %s\n",
			hours+"h",
			utils.ReplaceWeatherToIcons(entry.Weather[0].Description),
			int(entry.Main.Temp),
			entry.Main.Humidity,
			utils.HPaToMmHg(float64(entry.Main.Pressure)),
			entry.Wind.Speed,
			utils.DegreesToDirectionIcon(entry.Wind.Deg),
		)

		if hours == "21" || hours == "22" || hours == "23" {
			forecast += "\n"
		}

	}
	return forecast, nil
}