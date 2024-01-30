package service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/http_client"
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/service/convert"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	moreInfoURLFormat  = "\n\n<a href=\"https://openweathermap.org/city/%s\">🌐 More</a>"
	failedToGetWeather = "Failed to get weather data:"
	tryAnother         = "Please try another city name, or try sending the location."
)

type OWMService struct {
	httpClient http_client.HTTPClient
	cfg        *config.Config
	log        *logrus.Logger
}

func NewOpenWeatherMap(httpClient http_client.HTTPClient, cfg *config.Config, log *logrus.Logger) *OWMService {
	return &OWMService{
		httpClient: httpClient,
		cfg:        cfg,
		log:        log,
	}
}

// GetWeatherForecast retrieves the weather forecast based on the provided weather command,
// updates the user's last command, and returns the formatted weather message.
func (OW *OWMService) GetWeatherForecast(us model.UserData) (weatherMessage string, err error) {
	fc := "GetWeatherForecast"
	OW.log.Debug(us)

	var cityId string
	var weatherData model.WeatherCurrent
	var forecastData model.WeatherForecast

	getWeatherUrl, err := OW.generateWeatherUrl(us)
	if err != nil {
		OW.log.Errorf("%s: %v", fc, err)
		return "", err
	}

	r, err := OW.httpClient.Get(getWeatherUrl)
	if err != nil {
		OW.log.Errorf("%s: %v", fc, err)
		return "", err
	}

	if r.StatusCode != http.StatusOK {
		var errResponse model.ErrorResponse
		if r.StatusCode == http.StatusNotFound {
			decoder := json.NewDecoder(r.Body)
			if err = decoder.Decode(&errResponse); err != nil {
				OW.log.Errorf("%s: %v", fc, err)
				return "", err
			} else if err == nil {
				OW.log.Error(errResponse)
				return "", fmt.Errorf("%s", tryAnother)
			}
		}
		OW.log.Error(failedToGetWeather, r.StatusCode)
		return "", fmt.Errorf(failedToGetWeather, "%d", r.StatusCode)
	}

	decoder := json.NewDecoder(r.Body)
	if us.Last == model.CallbackCurrent || us.Last == model.CallbackCurrentLocation {
		if err = decoder.Decode(&weatherData); err != nil {
			OW.log.Errorf("%s: %v", fc, err)
			return "", err
		}
		weatherMessage, cityId = messageCurrentWeather(weatherData, us.Metric)
	} else if us.Last == model.CallbackForecast || us.Last == model.CallbackForecastLocation {
		if err = decoder.Decode(&forecastData); err != nil {
			OW.log.Errorf("%s: %v", fc, err)
			return "", err
		}
		weatherMessage, cityId = messageForecastWeather(forecastData, us.Metric)
	}

	return weatherMessage + fmt.Sprintf(moreInfoURLFormat, cityId), nil
}

// generateWeatherUrl ...
func (OW *OWMService) generateWeatherUrl(us model.UserData) (fullWeatherUrl string, err error) {
	fc := "generateWeatherUrl"
	OW.log.Debugf("%s command: %s", fc, us.Last)
	weatherUrl := OW.cfg.WeatherApiUrl
	if us.Last == model.CallbackCurrent || us.Last == model.CallbackCurrentLocation {
		weatherUrl += "weather?"
	} else if us.Last == model.CallbackForecast || us.Last == model.CallbackForecastLocation {
		weatherUrl += "forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		OW.log.Errorf("%s: can't parse the weather API URL, %v", fc, err)
		return "", err
	}

	q := url.Values{}

	if us.Last == model.CallbackCurrent || us.Last == model.CallbackForecast {
		q.Add("q", us.City)
	} else if us.Last == model.CallbackForecastLocation || us.Last == model.CallbackCurrentLocation {
		q.Add("lat", us.Lat)
		q.Add("lon", us.Lon)
	}

	q.Add("appid", OW.cfg.WToken)
	if us.Metric {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()

	OW.log.Debug(fc, " url.Values:", q)
	OW.log.Debug(fc, " weather resp.url:", u.String())

	return u.String(), nil
}

// messageCurrentWeather returns a message with current weather and city id (in string).
func messageCurrentWeather(current model.WeatherCurrent, metric bool) (messageWithCurrent, cityId string) {
	tUnits, wUnits, pUnits := convert.Units(metric)
	pressure := convert.Pressure(current.Main.Pressure, metric)
	windSpeed := convert.WindSpeed(current.Wind.Speed, metric)
	loc := time.FixedZone("Custom Timezone", current.Timezone)
	messageWithCurrent = fmt.Sprintf(
		"<b>%s %s</b> %s\n\n 🌡 %+d%s (Feel %+d%s) 💧 %d%%  \n\n 📉 %+d%s ️ 📈 %+d%s \n%d %s %.2f%s %s \n\n🌅  %s 🌉  %s",
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
	return messageWithCurrent, strconv.Itoa(current.ID)
}

// messageForecastWeather returns a message with weather forecast and city id (in string).
func messageForecastWeather(forecast model.WeatherForecast, metric bool) (message, cityIdStr string) {
	tUnits, wUnits, pUnits := convert.Units(metric)
	// A headers displaying the forecast country and city names,
	// along with units for time, temperature, humidity, pressure, and wind direction.
	headerPlace := fmt.Sprintf("<b>%s %s\n\n</b>", forecast.City.Country, forecast.City.Name)
	headerUnits := fmt.Sprintf("[TIME] [%s] [%s] [%s] [%s, dir.]\n", tUnits, "💧", pUnits, wUnits)
	message += headerPlace
	// Iterate through forecast entries to construct messages for each time period.
	for i, e := range forecast.List {
		forecastTime := time.Unix(int64(e.Dt), 0).
			In(time.FixedZone("Custom Timezone", forecast.City.Timezone))
		hours, day := forecastTime.Format("15"), forecastTime.Format("02")
		windSpeedForecast := convert.WindSpeed(e.Wind.Speed, metric)
		// Display the date for each day in the format:
		// 🗓 31 January (Wednesday) along with the header containing units.
		if hours == "01" || hours == "02" || i == 0 {
			message += fmt.Sprintf("<b>🗓 %s %s (%s)</b>\n",
				day, forecastTime.Month().String(), forecastTime.Weekday().String()) + headerUnits
		}
		// Hourly forecast.
		message += fmt.Sprintf("%s:00 %s %+d %d%% %d [%.1f %s]\n",
			hours,
			convert.AddIcon(e.Weather[0].Description, false),
			convert.KelvinToFahrenheitAndRound(e.Main.Temp, metric),
			e.Main.Humidity,
			convert.Pressure(e.Main.Pressure, metric),
			windSpeedForecast, convert.DegsToDirIcon(e.Wind.Deg))
		// Insert a newline between days at the end of each day.
		if hours == "21" || hours == "22" || hours == "23" {
			message += "\n"
		}
	}
	return message, strconv.Itoa(forecast.City.ID)
}
