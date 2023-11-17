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

var weatherData types.WeatherResponse
var forecastData types.WeatherResponse5d3h
var userMessage string

func GetWeather(fullUrlGet, forecastType string) (string, error) {
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
	if forecastType == "current" || forecastType == "current 📍" {
		err = json.Unmarshal(body, &weatherData)
		if err != nil {
			errorMessage := err.Error()
			log.Println("Error: ", errorMessage)
			return "", fmt.Errorf("error: %s", errorMessage)
		}
		userMessage = fmt.Sprintf("%s %s \n %s 🌡 %.0f°C 💧 %d%%\n\nFeel %.0f°C  📉 %.0f°C ️ 📈 %.0f°C \n %.2f mmHg %s %.2f m/s \n\n🌅  %s 🌉  %s",
			weatherData.Sys.Country,
			weatherData.Name,
			utils.ReplaceWeatherToIcons(weatherData.Weather[0].Description),
			weatherData.Main.Temp,
			weatherData.Main.Humidity,
			weatherData.Main.FeelsLike,
			weatherData.Main.TempMin,
			weatherData.Main.TempMax,
			utils.HPaToMmHg(float64(weatherData.Main.Pressure)),
			utils.DegreesToDirectionIcon(weatherData.Wind.Deg),
			weatherData.Wind.Speed,
			utils.TimeStampToHuman(weatherData.Sys.Sunrise, weatherData.Timezone, "15:04"),
			utils.TimeStampToHuman(weatherData.Sys.Sunset, weatherData.Timezone, "15:04"))

		cityId := strconv.Itoa(weatherData.ID)
		userMessage += "\n\n" + "More information on the web link:" + "\n" + "https://openweathermap.org/city/" + cityId

	} else if forecastType == "5-days forecast" || forecastType == "5-days forecast 📍" {
		err = json.Unmarshal(body, &forecastData)
		if err != nil {
			errorMessage := err.Error()
			log.Println("Error: ", errorMessage)
			return "", fmt.Errorf("error: %s", errorMessage)
		}
		// Creating a string to display the country and city names
		userMessage = fmt.Sprintf("<b>%s %s\n\n</b>", forecastData.City.Country, forecastData.City.Name)
		// Constructing the date display, including day, month, and day of the week,
		// to be inserted into the user message about the weather.
		userMessage += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", utils.TimeStampToHuman(forecastData.List[0].Dt, forecastData.City.Timezone, "02"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "m"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "d"))

		for ind, entry := range forecastData.List {
			hours := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "15")
			dayNum := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "02")
			dayOfWeek := utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "d")
			if hours == "01" || hours == "02" && ind > 0 {
				// Constructing the date display, including day, month, and day of the week,
				// to be inserted into the user message about the weather.
				userMessage += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", dayNum, utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "m"), dayOfWeek)
			}

			userMessage += fmt.Sprintf("%s %v %v°C %d%% %.1f mmHg %.1f m/s %s\n",
				hours,
				utils.ReplaceWeatherToIcons(entry.Weather[0].Description),
				int(entry.Main.Temp),
				entry.Main.Humidity,
				utils.HPaToMmHg(float64(entry.Main.Pressure)),
				entry.Wind.Speed,
				utils.DegreesToDirectionIcon(entry.Wind.Deg),
			)

			if hours == "21" || hours == "22" || hours == "23" {
				userMessage += "\n"
			}

		}
		cityId := strconv.Itoa(forecastData.City.ID)
		userMessage += "\n\n" + "More information on the web link:" + "\n" + "https://openweathermap.org/city/" + cityId
	}
	return userMessage, nil
}
