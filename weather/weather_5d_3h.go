package weather

import (
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func Get5DayForecast(fullUrlGet string) (string, error) {

	text9 := "More information on the web link:"

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
				return "", fmt.Errorf("%s. \nTry another city name.", errorResponse.Message)
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

	forecast += forecastData.City.Country + " " + forecastData.City.Name + "\n\n" + utils.TimeStampToHuman(forecastData.List[0].Dt, forecastData.City.Timezone, "02") + " " + utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "m") + " (" + utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "d") + ")\n"

	for _, entry := range forecastData.List {
		hours := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "15")
		dayNum := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "02")
		dayOfWeek := utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "d")
		if hours == "01" || hours == "02" {
			forecast += dayNum + " " + utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "m") + " (" + dayOfWeek + ")\n"
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
	cityId := strconv.Itoa(forecastData.City.ID)
	forecast += "\n\n" + text9 + "\n" + "https://openweathermap.org/city/" + cityId
	return forecast, nil
}
