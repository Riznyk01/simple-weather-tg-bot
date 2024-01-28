package weather_service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/weather_service/convert"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

// GetWeatherForecast ...
func (OW *OpenWeatherMapService) GetWeatherForecast(chatId int64, weatherCommand string) (weatherMessage string, err error) {
	fc := "GetWeatherForecast"

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

// messageCurrentWeather returns a message with current weather and city id (in string).
func messageCurrentWeather(current model.WeatherCurrent, metric bool) (userMessageCurrent, cityId string) {
	tUnits, wUnits, pUnits := convert.Units(metric)
	pressure := convert.Pressure(current.Main.Pressure, metric)
	windSpeed := convert.WindSpeed(current.Wind.Speed, metric)
	loc := time.FixedZone("Custom Timezone", current.Timezone)
	userMessageCurrent = fmt.Sprintf("<b>%s %s</b> %s\n\n 🌡 %+d%s (Feel %+d%s) 💧 %d%%  \n\n 📉 %+d%s ️ 📈 %+d%s \n%d %s %.2f%s %s \n\n🌅  %s 🌉  %s",
		current.Sys.Country, current.Name,
		convert.AddIcon(current.Weather[0].Description, true),
		convert.KelvinToFahrenheitAndRound(current.Main.Temp, metric), tUnits,
		convert.KelvinToFahrenheitAndRound(current.Main.FeelsLike, metric), tUnits,
		current.Main.Humidity,
		convert.KelvinToFahrenheitAndRound(current.Main.TempMin, metric), tUnits,
		convert.KelvinToFahrenheitAndRound(current.Main.TempMax, metric), tUnits,
		pressure, pUnits,
		windSpeed, wUnits, convert.DegsToDirIcon(current.Wind.Deg),
		time.Unix(int64(current.Sys.Sunrise), 0).In(loc).Format("15:04"),
		time.Unix(int64(current.Sys.Sunset), 0).In(loc).Format("15:04"))
	return userMessageCurrent, strconv.Itoa(current.ID)
}

// messageForecastWeather returns a message with weather forecast and city id (in string).
func messageForecastWeather(forecast model.WeatherForecast, metric bool) (message, cityIdStr string) {
	tUnits, wUnits, pUnits := convert.Units(metric)
	// A headers displaying the forecast country and city names,
	// along with units for time, temperature, humidity, pressure, and wind direction.
	headerPlace := fmt.Sprintf("<b>%s %s\n\n</b>", forecast.City.Country, forecast.City.Name)
	headerUnits := fmt.Sprintf("[HOURS] [%s] [%s] [%s] [%s, dir.]\n", tUnits, "💧", pUnits, wUnits)
	message += headerPlace
	// ...
	for ind, entry := range forecast.List {
		forecastTime := time.Unix(int64(entry.Dt), 0).
			In(time.FixedZone("Custom Timezone", forecast.City.Timezone))
		hours, day := forecastTime.Format("15"), forecastTime.Format("02")
		windSpeedForecast := convert.WindSpeed(entry.Wind.Speed, metric)

		if hours == "01" || hours == "02" || ind == 0 {
			// The date for each day is displayed in the format: 🗓 31 January (Wednesday).
			message += fmt.Sprintf("<b>🗓 %s %s (%s)</b>\n",
				day, forecastTime.Month().String(), forecastTime.Weekday().String())
			message += headerUnits
		}

		message += fmt.Sprintf("%s:00 %s %+d %d%% %d [%.1f %s]\n",
			hours,
			convert.AddIcon(entry.Weather[0].Description, false),
			convert.KelvinToFahrenheitAndRound(entry.Main.Temp, metric),
			entry.Main.Humidity,
			convert.Pressure(entry.Main.Pressure, metric),
			windSpeedForecast, convert.DegsToDirIcon(entry.Wind.Deg),
		)

		if hours == "21" || hours == "22" || hours == "23" {
			message += "\n"
		}

	}
	return message, strconv.Itoa(forecast.City.ID)
}
