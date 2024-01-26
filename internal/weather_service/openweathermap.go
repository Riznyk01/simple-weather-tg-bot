package weather_service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/weather_service/util"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	moreInfoURLFormat  = "\n\n<a href=\"https://openweathermap.org/city/%s\">🌐 More</a>"
	failedToGetWeather = "Failed to get weather data:"
	tryAnother         = "Please try another city name, or try sending the location."
	systemFetchError   = "Encountered an error when trying to fetch the users system of measurement:"
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
	return OW.repo.SetSystem(chatId, system)
}
func (OW *OpenWeatherMapService) SetCity(chatId int64, city string) error {
	return OW.repo.SetCity(chatId, city)
}
func (OW *OpenWeatherMapService) SetLocation(chatId int64, lat, lon string) error {
	return OW.repo.SetLocation(chatId, lat, lon)
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

// SetLast returns a complete weather message and sets last weather responce.
func (OW *OpenWeatherMapService) SetLast(chatId int64, weatherCommand string) (weatherMessage string, err error) {
	fc := "SetLast"

	weatherUrl := OW.cfg.WeatherApiUrl
	var cityId string
	var weatherData model.WeatherCurrent
	var forecastData model.WeatherForecast
	var errResponse model.ErrorResponse

	switch weatherCommand {
	case model.CallbackCurrent, model.CallbackCurrentLocation:
		weatherUrl += "weather?"
	case model.CallbackForecast, model.CallbackForecastLocation:
		weatherUrl += "forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		return "", err
	}

	q := url.Values{}

	switch weatherCommand {
	case model.CallbackCurrent, model.CallbackForecast:
		city, err := OW.repo.GetCity(chatId)
		if err != nil {
			return "", err
		}
		q.Add("q", city)
	case model.CallbackForecastLocation, model.CallbackCurrentLocation:
		lat, lon, err := OW.repo.GetLocation(chatId)
		if err != nil {
			return "", err
		}
		q.Add("lat", lat)
		q.Add("lon", lon)
	}

	q.Add("appid", OW.cfg.WToken)
	metric, err := OW.repo.GetSystem(chatId)
	if err != nil {
		OW.log.Error(systemFetchError, err)
	}
	if metric {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()

	OW.log.Debug(fc, " url.Values:", q)
	OW.log.Debug(fc, " weather resp.url:", u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		OW.log.Error(fc, ": ", err)
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
			err = json.Unmarshal(body, &errResponse)
			if err == nil {
				OW.log.Error(errResponse)
				return "", fmt.Errorf("%s", tryAnother)
			}
		}
		OW.log.Error(failedToGetWeather, resp.StatusCode)
		return "", fmt.Errorf(failedToGetWeather, "%d", resp.StatusCode)
	}

	switch weatherCommand {
	case model.CallbackCurrent, model.CallbackCurrentLocation:
		err = json.Unmarshal(body, &weatherData)
		if err != nil {
			OW.log.Error(err)
			return "", err
		}
		weatherMessage, cityId = messageCurrentWeather(weatherData, metric)
	case model.CallbackForecast, model.CallbackForecastLocation:
		err = json.Unmarshal(body, &forecastData)
		if err != nil {
			OW.log.Error(err)
			return "", err
		}
		weatherMessage, cityId = messageForecastWeather(forecastData, metric)
	}

	err = OW.repo.SetLast(chatId, weatherCommand)
	if err != nil {
		return "", err
	}
	return weatherMessage + fmt.Sprintf(moreInfoURLFormat, cityId), nil
}
func (OW *OpenWeatherMapService) GetLast(chatId int64) (weatherMessage string, err error) {
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

// units returns units based on the metric system.
func units(metricUnits bool) (tempUnits, windUnits, pressureUnits string) {
	if metricUnits {
		return "°C", "m/s", "mmHg"
	}
	return "°F", "mph", "inHg"
}

// messageCurrentWeather returns a message with current weather and city id (in string).
func messageCurrentWeather(currentData model.WeatherCurrent, metric bool) (userMessageCurrent, cityId string) {
	tUnits, wUnits, pUnits := units(metric)
	pressure := util.PressureConverting(currentData.Main.Pressure, metric)
	windSpeed := currentData.Wind.Speed
	//Converting to miles per hour if non-metric
	if !metric {
		windSpeed = util.ToMilesPerHour(currentData.Wind.Speed)
		pressure = util.PressureConverting(currentData.Main.Pressure, metric)
	}
	userMessageCurrent = fmt.Sprintf("<b>%s %s</b> %s\n\n 🌡 %+d%s (Feel %+d%s) 💧 %d%%  \n\n 📉 %+d%s ️ 📈 %+d%s \n%d %s %.2f%s %s \n\n🌅  %s 🌉  %s",
		currentData.Sys.Country,
		currentData.Name,
		util.WeatherTextToIcon(currentData.Weather[0].Description, true),
		util.TemperatureConverting(currentData.Main.Temp, metric),
		tUnits,
		util.TemperatureConverting(currentData.Main.FeelsLike, metric),
		tUnits,
		currentData.Main.Humidity,
		util.TemperatureConverting(currentData.Main.TempMin, metric),
		tUnits,
		util.TemperatureConverting(currentData.Main.TempMax, metric),
		tUnits,
		pressure,
		pUnits,
		windSpeed,
		wUnits,
		util.DegreesToDirectionIcon(currentData.Wind.Deg),
		util.TimeStampToHuman(currentData.Sys.Sunrise, currentData.Timezone, "15:04"),
		util.TimeStampToHuman(currentData.Sys.Sunset, currentData.Timezone, "15:04"))

	return userMessageCurrent, strconv.Itoa(currentData.ID)
}

// messageForecastWeather returns a message with weather forecast and city id (in string).
func messageForecastWeather(forecastData model.WeatherForecast, metric bool) (message, cityIdStr string) {
	tUnits, wUnits, pUnits := units(metric)
	// Creating a string to display the country and city names
	message = fmt.Sprintf("<b>%s %s\n\n</b>", forecastData.City.Country, forecastData.City.Name)
	// Constructing the date display, including day, month, and day of the week,
	// to be inserted into the user message about the weather.
	message += fmt.Sprintf("<b>🗓 %s %s (%s)</b>\n",
		util.TimeStampToHuman(forecastData.List[0].Dt, forecastData.City.Timezone, "02"),
		util.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "m"),
		util.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "d"))
	messageHeader := fmt.Sprintf("[time] [   ] [%s] [%s] [%s] [%s, dir.]\n",
		tUnits, "💧", pUnits, wUnits)
	message += messageHeader

	for ind, entry := range forecastData.List {
		hours := util.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "15")
		dayNum := util.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "02")
		dayOfWeek := util.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "d")
		if hours == "01" || hours == "02" && ind > 0 {
			// Constructing the date display, including day, month, and day of the week,
			// to be inserted into the user message about the weather.
			message += fmt.Sprintf("<b>🗓 %s %s (%s)</b>\n", dayNum, util.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "m"), dayOfWeek)
			message += messageHeader
		}

		windSpeedForecast := entry.Wind.Speed
		// Converting to miles per hour if non-metric
		if !metric {
			windSpeedForecast = util.ToMilesPerHour(entry.Wind.Speed)
		}

		message += fmt.Sprintf("%s %s %+d %d%% %d [%.1f %s]\n",
			hours+":00",
			util.WeatherTextToIcon(entry.Weather[0].Description, false),
			util.TemperatureConverting(entry.Main.Temp, metric),
			entry.Main.Humidity,
			util.PressureConverting(entry.Main.Pressure, metric),
			windSpeedForecast,
			util.DegreesToDirectionIcon(entry.Wind.Deg),
		)

		if hours == "21" || hours == "22" || hours == "23" {
			message += "\n"
		}

	}
	return message, strconv.Itoa(forecastData.City.ID)
}
