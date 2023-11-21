package weather

import (
	"SimpleWeatherTgBot/lib/e"
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var weatherData types.WeatherResponse
var forecastData types.WeatherResponse5d3h
var userMessage, temperatureUnits, windUnits, weatherMessage string

func GetWeather(fullUrlGet, forecastType string, metric bool) (string, error) {
	resp, err := http.Get(fullUrlGet)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := err.Error()
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
				return "", fmt.Errorf("%s. Try another city name.", errorResponse.Message)
			}
		}
		return "", fmt.Errorf("Failed to get weather data. Status code: %d", resp.StatusCode)
	}
	//Getting units of measurement
	temperatureUnits, windUnits = units(metric)

	if forecastType == "current" || forecastType == "current 📍" {
		err = json.Unmarshal(body, &weatherData)
		if err != nil {
			return "", e.Wrap("", err)
		}
		weatherMessage = messageCurrentWeather(weatherData, metric)
	} else if forecastType == "5-days forecast" || forecastType == "5-days forecast 📍" {
		err = json.Unmarshal(body, &forecastData)
		if err != nil {
			return "", e.Wrap("", err)
		}
		weatherMessage = messageForecastWeather(forecastData, metric)
	}
	return weatherMessage, nil
}

// Returns the units depending on whether the system is metric or not
func units(metricUnits bool) (tempUnits, windUnits string) {
	if metricUnits {
		tempUnits = " °C"
		windUnits = " m/s"
	} else {
		tempUnits = " °F"
		windUnits = " mph"
	}
	return tempUnits, windUnits
}
func messageCurrentWeather(weatherData types.WeatherResponse, metric bool) string {
	//Convetring to miles per hour if non metric
	windSpeed := weatherData.Wind.Speed
	if !metric {
		windSpeed = utils.ToMilesPerHour(weatherData.Wind.Speed)
	}

	userMessageCurrent := fmt.Sprintf("<b>%s %s</b> %s 🌡 %.0f%s (Feel %.0f%s) 💧 %d%%  \n\n 📉 %.0f%s ️ 📈 %.0f%s %.2f mmHg %.2f%s %s \n\n🌅  %s 🌉  %s",
		weatherData.Sys.Country,
		weatherData.Name,
		utils.ReplaceWeatherToIcons(weatherData.Weather[0].Description),
		weatherData.Main.Temp,
		temperatureUnits,
		weatherData.Main.FeelsLike,
		temperatureUnits,
		weatherData.Main.Humidity,
		weatherData.Main.TempMin,
		temperatureUnits,
		weatherData.Main.TempMax,
		temperatureUnits,
		utils.HPaToMmHg(float64(weatherData.Main.Pressure)),
		windSpeed,
		windUnits,
		utils.DegreesToDirectionIcon(weatherData.Wind.Deg),
		utils.TimeStampToHuman(weatherData.Sys.Sunrise, weatherData.Timezone, "15:04"),
		utils.TimeStampToHuman(weatherData.Sys.Sunset, weatherData.Timezone, "15:04"))

	cityId := strconv.Itoa(weatherData.ID)
	userMessageCurrent += " More: " + "<a href=\"https://openweathermap.org/city/" + cityId + "\">OpenWeatherMap</a>"

	return userMessageCurrent
}
func messageForecastWeather(forecastData types.WeatherResponse5d3h, metric bool) string {
	// Creating a string to display the country and city names
	userMessageForecast := fmt.Sprintf("<b>%s %s\n\n</b>", forecastData.City.Country, forecastData.City.Name)
	// Constructing the date display, including day, month, and day of the week,
	// to be inserted into the user message about the weather.
	userMessageForecast += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", utils.TimeStampToHuman(forecastData.List[0].Dt, forecastData.City.Timezone, "02"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "m"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "d"))
	messageHeader := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s]\n",
		"h:m",
		temperatureUnits,
		"%",
		"mmHg",
		windUnits,
	)

	userMessageForecast += messageHeader

	for ind, entry := range forecastData.List {
		hours := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "15")
		dayNum := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "02")
		dayOfWeek := utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "d")
		if hours == "01" || hours == "02" && ind > 0 {
			// Constructing the date display, including day, month, and day of the week,
			// to be inserted into the user message about the weather.
			userMessageForecast += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", dayNum, utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "m"), dayOfWeek)
			userMessageForecast += messageHeader
		}

		windSpeedForecast := entry.Wind.Speed
		//Convetring to miles per hour if non metric
		if !metric {
			windSpeedForecast = utils.ToMilesPerHour(entry.Wind.Speed)
		}

		userMessageForecast += fmt.Sprintf("%s %s %+d %d%% %.1f %.1f %s\n",
			hours+":00",
			utils.ReplaceWeatherToIcons(entry.Weather[0].Description),
			int(entry.Main.Temp),
			entry.Main.Humidity,
			utils.HPaToMmHg(float64(entry.Main.Pressure)),
			windSpeedForecast,
			utils.DegreesToDirectionIcon(entry.Wind.Deg),
		)

		if hours == "21" || hours == "22" || hours == "23" {
			userMessageForecast += "\n"
		}

	}
	cityId := strconv.Itoa(forecastData.City.ID)
	userMessageForecast += "\nMore: " + "<a href=\"https://openweathermap.org/city/" + cityId + "\">OpenWeatherMap</a>"
	return userMessageForecast
}
