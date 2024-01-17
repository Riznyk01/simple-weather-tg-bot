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

const (
	moreInfoURLFormat = "\n\n<a href=\"https://openweathermap.org/city/%s\">🌐 More</a>"
	failed            = "Failed to get weather data. Status code:"
	tryAnother        = "Please try another city name, or try sending the location."
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

func (OW *OpenWeatherMapService) SetSystem(chatId int64, system bool) error {
	if err := OW.repo.SetSystem(chatId, system); err != nil {
		return err
	}
	return nil
}
func (OW *OpenWeatherMapService) SetCity(chatId int64, city string) error {
	if err := OW.repo.SetCity(chatId, city); err != nil {
		return err
	}
	return nil
}
func (OW *OpenWeatherMapService) SetLocation(chatId int64, lat, lon string) error {
	if err := OW.repo.SetLocation(chatId, lat, lon); err != nil {
		return err
	}
	return nil
}
func (OW *OpenWeatherMapService) GetSystem(chatId int64) (bool, error) {
	return OW.repo.GetSystem(chatId)
}
func (OW *OpenWeatherMapService) GetCity(chatId int64) (string, error) {
	return OW.repo.GetCity(chatId)
}
func (OW *OpenWeatherMapService) GetLocation(chatId int64) (string, string, error) {
	return OW.repo.GetLocation(chatId)
}
func (OW *OpenWeatherMapService) Exists(chatId int64) (bool, error) {
	return OW.repo.Exists(chatId)
}

// Returns a complete weather message and sets last weather responce.
func (OW *OpenWeatherMapService) SetLast(chatId int64, weatherCommand string) (weatherMessage string, err error) {
	weatherUrl := OW.cfg.WeatherApiUrl
	var cityIdString string
	var weatherData types.WeatherCurrent
	var forecastData types.WeatherForecast
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
		city, err := OW.repo.GetCity(chatId)
		if err != nil {
			return "", err
		}
		if city == "" {
			return "empty", err
		}
		q.Add("q", city)
	} else if weatherCommand == types.CommandForecastLocation || weatherCommand == types.CommandCurrentLocation {
		lat, lon, err := OW.repo.GetLocation(chatId)
		if err != nil {
			return "", err
		}
		if lat == "" || lon == "" {
			return "empty", err
		}
		q.Add("lat", lat)
		q.Add("lon", lon)
	}

	q.Add("appid", OW.cfg.WToken)
	metric, err := OW.repo.GetSystem(chatId)
	if err != nil {
		OW.log.Info("Encountered an error when trying to fetch the users system of measurement:", err)
	}
	if metric {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()

	OW.log.Debug("url.Values:", q)
	OW.log.Debug("weather response url:", u.String())
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
				return "", fmt.Errorf("%s", tryAnother)
			}
		}
		return "", fmt.Errorf(failed, "%d", resp.StatusCode)
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
	more := fmt.Sprintf(moreInfoURLFormat, cityIdString)
	err = OW.repo.SetLast(chatId, weatherCommand)
	if err != nil {
		return "", err
	}
	return weatherMessage + more, nil
}
func (OW *OpenWeatherMapService) GetLast(chatId int64) (weatherMessage string, err error) {
	ex, err := OW.repo.Exists(chatId)
	if err != nil {
		return "", err
	}
	if !ex {
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
func (OW *OpenWeatherMapService) AddRequestsCount(chatId int64) int {
	return OW.repo.AddRequestsCount(chatId)
}

// Returns units based on the metric system.
func units(metricUnits bool) (tempUnits, windUnits, pressureUnits string) {
	if metricUnits {
		return "°C", "m/s", "mmHg"
	}
	return "°F", "mph", "inHg"
}

// Returns a message with current weather and city id (in string).
func messageCurrentWeather(weatherData types.WeatherCurrent, metric bool) (userMessageCurrent, cityIdStr string) {
	tUnits, wUnits, pUnits := units(metric)
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
		tUnits,
		utils.TemperatureConverting(weatherData.Main.FeelsLike, metric),
		tUnits,
		weatherData.Main.Humidity,
		utils.TemperatureConverting(weatherData.Main.TempMin, metric),
		tUnits,
		utils.TemperatureConverting(weatherData.Main.TempMax, metric),
		tUnits,
		pressure,
		pUnits,
		windSpeed,
		wUnits,
		utils.DegreesToDirectionIcon(weatherData.Wind.Deg),
		utils.TimeStampToHuman(weatherData.Sys.Sunrise, weatherData.Timezone, "15:04"),
		utils.TimeStampToHuman(weatherData.Sys.Sunset, weatherData.Timezone, "15:04"))

	return userMessageCurrent, strconv.Itoa(weatherData.ID)
}

// Returns a message with weather forecast and city id (in string).
func messageForecastWeather(forecastData types.WeatherForecast, metric bool) (message, cityIdStr string) {
	tUnits, wUnits, pUnits := units(metric)
	// Creating a string to display the country and city names
	message = fmt.Sprintf("<b>%s %s\n\n</b>", forecastData.City.Country, forecastData.City.Name)
	// Constructing the date display, including day, month, and day of the week,
	// to be inserted into the user message about the weather.
	message += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", utils.TimeStampToHuman(forecastData.List[0].Dt, forecastData.City.Timezone, "02"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "m"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "d"))
	messageHeader := fmt.Sprintf("[h:m] [---] [%s] [%s] [%s] [%s]\n",
		tUnits, "%", pUnits, wUnits)
	message += messageHeader

	for ind, entry := range forecastData.List {
		hours := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "15")
		dayNum := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "02")
		dayOfWeek := utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "d")
		if hours == "01" || hours == "02" && ind > 0 {
			// Constructing the date display, including day, month, and day of the week,
			// to be inserted into the user message about the weather.
			message += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", dayNum, utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "m"), dayOfWeek)
			message += messageHeader
		}

		windSpeedForecast := entry.Wind.Speed
		// Converting to miles per hour if non-metric
		if !metric {
			windSpeedForecast = utils.ToMilesPerHour(entry.Wind.Speed)
		}

		message += fmt.Sprintf("%s %s %+d %d%% %.1f %.1f %s\n",
			hours+":00",
			utils.ReplaceWeatherToIcons(entry.Weather[0].Description),
			utils.TemperatureConverting(entry.Main.Temp, metric),
			entry.Main.Humidity,
			utils.PressureConverting(float64(entry.Main.Pressure), metric),
			windSpeedForecast,
			utils.DegreesToDirectionIcon(entry.Wind.Deg),
		)

		if hours == "21" || hours == "22" || hours == "23" {
			message += "\n"
		}

	}
	return message, strconv.Itoa(forecastData.City.ID)
}
