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

func GetWeather(fullUrlGet string) (string, error) {

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

	userMessage := fmt.Sprintf("%s %s - %s 🌡 %.1f°C 💧 %d%%\n\nFeel %.1f°C  📉 %.1f°C ️ 📈 %.1f°C \n %.2f mmHg %s %.2f m/s \n\n🌅  %s 🌉  %s",
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
	userMessage += "\n\n" + "https://openweathermap.org/city/" + cityId

	return userMessage, nil
}
