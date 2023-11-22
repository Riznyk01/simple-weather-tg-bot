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
var temperatureUnits, windUnits, pressureUnits string

// Returns a complete weather message.
func GetWeather(fullUrlGet, forecastType string, metric bool) (weatherMessage string, err error) {
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

	var cityIdString string

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
	temperatureUnits, windUnits, pressureUnits = units(metric)
	if forecastType == "current" || forecastType == "current 📍" {
		err = json.Unmarshal(body, &weatherData)
		if err != nil {
			return "", e.Wrap("", err)
		}
		weatherMessage, cityIdString = messageCurrentWeather(weatherData, metric)
	} else if forecastType == "5-days forecast" || forecastType == "5-days forecast 📍" {
		err = json.Unmarshal(body, &forecastData)
		if err != nil {
			return "", e.Wrap("", err)
		}
		weatherMessage, cityIdString = messageForecastWeather(forecastData, metric)
	}
	more := fmt.Sprintf("\nMore: <a href=\"https://openweathermap.org/city/%s\">OpenWeatherMap</a>", cityIdString)
	return weatherMessage + more, nil
}

// Returns units based on the metric system.
func units(metricUnits bool) (tempUnits, windUnits, pressureUnits string) {
	if metricUnits {
		tempUnits = " °C"
		windUnits = " m/s"
		pressureUnits = " mmHg"
	} else {
		tempUnits = " °F"
		windUnits = " mph"
		pressureUnits = " inHg"
	}
	return tempUnits, windUnits, pressureUnits
}

// Returns a message with current weather and city id (in string).
func messageCurrentWeather(weatherData types.WeatherResponse, metric bool) (userMessageCurrent, cityIdStr string) {
	pressure := utils.PressureConverting(float64(weatherData.Main.Pressure), metric)
	windSpeed := weatherData.Wind.Speed
	//Converting to miles per hour if non-metric
	if !metric {
		windSpeed = utils.ToMilesPerHour(weatherData.Wind.Speed)
		pressure = utils.PressureConverting(float64(weatherData.Main.Pressure), metric)
	}
	userMessageCurrent = fmt.Sprintf("<b>%s %s</b> %s\n\n 🌡 %.0f%s (Feel %.0f%s) 💧 %d%%  \n\n 📉 %.0f%s ️ 📈 %.0f%s \n%.0f %s %.2f%s %s \n\n🌅  %s 🌉  %s",
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
		pressure,
		pressureUnits,
		windSpeed,
		windUnits,
		utils.DegreesToDirectionIcon(weatherData.Wind.Deg),
		utils.TimeStampToHuman(weatherData.Sys.Sunrise, weatherData.Timezone, "15:04"),
		utils.TimeStampToHuman(weatherData.Sys.Sunset, weatherData.Timezone, "15:04"))

	return userMessageCurrent, strconv.Itoa(weatherData.ID)
}

// Returns a message with weather forecast and city id (in string).
func messageForecastWeather(forecastData types.WeatherResponse5d3h, metric bool) (userMessageForecast, cityIdStr string) {
	// Creating a string to display the country and city names
	userMessageForecast = fmt.Sprintf("<b>%s %s\n\n</b>", forecastData.City.Country, forecastData.City.Name)
	// Constructing the date display, including day, month, and day of the week,
	// to be inserted into the user message about the weather.
	userMessageForecast += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", utils.TimeStampToHuman(forecastData.List[0].Dt, forecastData.City.Timezone, "02"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "m"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "d"))
	messageHeader := fmt.Sprintf("[%s] [---] [%s] [%s] [%s] [%s]\n",
		"h:m",
		temperatureUnits,
		"%",
		pressureUnits,
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

		pressure := utils.PressureConverting(float64(entry.Main.Pressure), metric)
		windSpeedForecast := entry.Wind.Speed
		// Converting to miles per hour if non-metric
		if !metric {
			pressure = utils.PressureConverting(float64(entry.Main.Pressure), metric)
			windSpeedForecast = utils.ToMilesPerHour(entry.Wind.Speed)
		}

		userMessageForecast += fmt.Sprintf("%s %s %+d %d%% %.1f %.1f %s\n",
			hours+":00",
			utils.ReplaceWeatherToIcons(entry.Weather[0].Description),
			int(entry.Main.Temp),
			entry.Main.Humidity,
			pressure,
			windSpeedForecast,
			utils.DegreesToDirectionIcon(entry.Wind.Deg),
		)

		if hours == "21" || hours == "22" || hours == "23" {
			userMessageForecast += "\n"
		}

	}
	return userMessageForecast, strconv.Itoa(forecastData.City.ID)
}
