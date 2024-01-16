package weather_service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/repository"
	"SimpleWeatherTgBot/types"
	"SimpleWeatherTgBot/utils"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type OpenWeatherMapService struct {
	repo *repository.Repository
	cfg  *config.Config
	log  *logrus.Logger
}

func NewOpenWeatherMap(repo *repository.Repository, cfg *config.Config, log *logrus.Logger) *OpenWeatherMapService {
	return &OpenWeatherMapService{
		repo: repo,
		cfg:  cfg,
		log:  log,
	}
}

func (OW *OpenWeatherMapService) SetSystem(chatId int64, system bool) {
	OW.repo.SetSystem(chatId, system)
}
func (OW *OpenWeatherMapService) SetCity(chatId int64, city string) {
	OW.repo.SetCity(chatId, city)
}
func (OW *OpenWeatherMapService) SetLocation(chatId int64, lat, lon string) {
	OW.repo.SetLocation(chatId, lat, lon)
}
func (OW *OpenWeatherMapService) GetSystem(chatId int64) (bool, error) {
	return OW.repo.GetSystem(chatId)
}
func (OW *OpenWeatherMapService) GetCity(chatId int64) string {
	return OW.repo.GetCity(chatId)
}
func (OW *OpenWeatherMapService) GetLat(chatId int64) string {
	return OW.repo.GetLat(chatId)
}
func (OW *OpenWeatherMapService) GetLon(chatId int64) string {
	return OW.repo.GetLon(chatId)
}
func (OW *OpenWeatherMapService) Exists(chatId int64) bool {
	return OW.repo.Exists(chatId)
}

// Returns a complete weather message and sets last weather responce.
func (OW *OpenWeatherMapService) SetLast(chatId int64, weatherCommand string) (weatherMessage string, err error) {
	weatherUrl := OW.cfg.WeatherApiUrl
	var cityIdString string
	var weatherData types.WeatherResponse
	var forecastData types.WeatherResponse5d3h
	if weatherCommand == types.CommandCurrent || weatherCommand == types.CommandCurrentLocation {
		weatherUrl += "weather?"
	} else if weatherCommand == types.CommandForecast || weatherCommand == types.CommandForecastLocation {
		weatherUrl += "forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		return "", err
	}

	q := url.Values{}

	if weatherCommand == types.CommandCurrent || weatherCommand == types.CommandForecast {
		city := OW.repo.GetCity(chatId)
		if city == "" {
			return "empty", err
		}
		q.Add("q", city)
	} else if weatherCommand == types.CommandForecastLocation || weatherCommand == types.CommandCurrentLocation {
		lat, lon := OW.repo.GetLat(chatId), OW.repo.GetLon(chatId)
		if lat == "" || lon == "" {
			return "empty", err
		}
		q.Add("lat", lat)
		q.Add("lon", lon)
	}

	q.Add("appid", OW.cfg.WToken)
	metric, err := OW.repo.GetSystem(chatId)
	if err != nil {
		//log.info
	}
	if metric {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()

	OW.log.Info(q)
	OW.log.Info(u.String())
	resp, err := http.Get(u.String())
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
	if weatherCommand == types.CommandCurrent || weatherCommand == types.CommandCurrentLocation {
		err = json.Unmarshal(body, &weatherData)
		if err != nil {
			return "", err
		}
		weatherMessage, cityIdString = messageCurrentWeather(weatherData, metric)
	} else if weatherCommand == types.CommandForecast || weatherCommand == types.CommandForecastLocation {
		err = json.Unmarshal(body, &forecastData)
		if err != nil {
			return "", err
		}
		weatherMessage, cityIdString = messageForecastWeather(forecastData, metric)
	}
	more := fmt.Sprintf("\n\n<a href=\"https://openweathermap.org/city/%s\">🌐 More</a>", cityIdString)
	OW.repo.SetLast(chatId, weatherCommand)
	return weatherMessage + more, nil
}
func (OW *OpenWeatherMapService) GetLast(chatId int64) (weatherMessage string, err error) {
	if !OW.repo.Exists(chatId) {
		return "empty", err
	}
	weatherCommand, err := OW.repo.GetLast(chatId)
	if err != nil {
		return "", err
	}
	weatherMessage, err = OW.SetLast(chatId, weatherCommand)
	if err != nil {
		return weatherMessage, err
	}
	return weatherMessage, nil
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
	temperatureUnits, windUnits, pressureUnits := units(metric)
	pressure := utils.PressureConverting(float64(weatherData.Main.Pressure), metric)
	windSpeed := weatherData.Wind.Speed
	//Converting to miles per hour if non-metric
	if !metric {
		windSpeed = utils.ToMilesPerHour(weatherData.Wind.Speed)
		pressure = utils.PressureConverting(float64(weatherData.Main.Pressure), metric)
	}
	userMessageCurrent = fmt.Sprintf("<b>%s %s</b> %s\n\n 🌡 %+d%s (Feel %+d%s) 💧 %d%%  \n\n 📉 %+d%s ️ 📈 %+d%s \n%.0f %s %.2f%s %s \n\n🌅  %s 🌉  %s",
		weatherData.Sys.Country,
		weatherData.Name,
		utils.ReplaceWeatherPlusIcons(weatherData.Weather[0].Description),
		utils.TemperatureConverting(weatherData.Main.Temp, metric),
		temperatureUnits,
		utils.TemperatureConverting(weatherData.Main.FeelsLike, metric),
		temperatureUnits,
		weatherData.Main.Humidity,
		utils.TemperatureConverting(weatherData.Main.TempMin, metric),
		temperatureUnits,
		utils.TemperatureConverting(weatherData.Main.TempMax, metric),
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
	temperatureUnits, windUnits, pressureUnits := units(metric)
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

		windSpeedForecast := entry.Wind.Speed
		// Converting to miles per hour if non-metric
		if !metric {
			windSpeedForecast = utils.ToMilesPerHour(entry.Wind.Speed)
		}

		userMessageForecast += fmt.Sprintf("%s %s %+d %d%% %.1f %.1f %s\n",
			hours+":00",
			utils.ReplaceWeatherToIcons(entry.Weather[0].Description),
			utils.TemperatureConverting(entry.Main.Temp, metric),
			entry.Main.Humidity,
			utils.PressureConverting(float64(entry.Main.Pressure), metric),
			windSpeedForecast,
			utils.DegreesToDirectionIcon(entry.Wind.Deg),
		)

		if hours == "21" || hours == "22" || hours == "23" {
			userMessageForecast += "\n"
		}

	}
	return userMessageForecast, strconv.Itoa(forecastData.City.ID)
}
