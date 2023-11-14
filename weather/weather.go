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

func GetWeather(city, tWeather string) (string, error) {

	weatherUrl := "https://api.openweathermap.org/data/2.5/weather?"

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

	var weatherData types.WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		return "", fmt.Errorf("error: %s", errorMessage)
	}

	userMessage := fmt.Sprintf("%s %s - %s 🌡 %.1f°C 💧 %d%%\n\nFeel %.1f°C  📉 %.1f°C ️ 📈 %.1f°C \n %.2f mmHg %.2f m/s (%s) \n\n🌅  %s 🌉  %s",
		weatherData.Sys.Country,
		weatherData.Name,
		utils.AddWeatherIcons(weatherData.Weather[0].Main),
		weatherData.Main.Temp,
		weatherData.Main.Humidity,
		weatherData.Main.FeelsLike,
		weatherData.Main.TempMin,
		weatherData.Main.TempMax,
		utils.HPaToMmHg(float64(weatherData.Main.Pressure)),
		weatherData.Wind.Speed,
		utils.DegreesToDirection(weatherData.Wind.Deg),
		utils.TimeStampToHuman(weatherData.Sys.Sunrise, weatherData.Timezone, "15:04"),
		utils.TimeStampToHuman(weatherData.Sys.Sunset, weatherData.Timezone, "15:04"))

	return userMessage, nil
}
