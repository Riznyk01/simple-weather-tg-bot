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

// GetWeatherForecast retrieves the weather forecast based on the provided weather command,
// updates the user's last command, and returns the formatted weather message.
func (OW *OpenWeatherMapService) GetWeatherForecast(chatId int64, command string) (weatherMessage string, err error) {

	var cityId string
	var weatherData model.WeatherCurrent
	var forecastData model.WeatherForecast

	metric, err := OW.repo.GetSystem(chatId)
	if err != nil {
		OW.log.Error(systemFetchError, err)
	}

	getWeatherUrl, err := OW.generateWeatherUrl(chatId, command)
	if err != nil {
		return "", err
	}

	respBody, err := OW.makeHttpRequest(getWeatherUrl)
	if err != nil {
		return "", err
	}

	switch command {
	case model.CallbackCurrent, model.CallbackCurrentLocation:
		err = json.Unmarshal(respBody, &weatherData)
		if err != nil {
			OW.log.Error(err)
			return "", err
		}
		weatherMessage, cityId = messageCurrentWeather(weatherData, metric)
	case model.CallbackForecast, model.CallbackForecastLocation:
		err = json.Unmarshal(respBody, &forecastData)
		if err != nil {
			OW.log.Error(err)
			return "", err
		}
		weatherMessage, cityId = messageForecastWeather(forecastData, metric)
	}

	err = OW.repo.SetLast(chatId, command)
	if err != nil {
		return "", err
	}
	return weatherMessage + fmt.Sprintf(moreInfoURLFormat, cityId), nil
}

// generateWeatherUrl ...
func (OW *OpenWeatherMapService) generateWeatherUrl(chatId int64, command string) (fullWeatherUrl string, err error) {
	fc := "generateWeatherUrl"

	weatherUrl := OW.cfg.WeatherApiUrl
	switch command {
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

	switch command {
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

	fullWeatherUrl = u.String()

	return fullWeatherUrl, nil
}

// makeHttpRequest ...
func (OW *OpenWeatherMapService) makeHttpRequest(fullUrl string) ([]byte, error) {
	fc := "makeHttpRequest"

	var errResponse model.ErrorResponse

	resp, err := http.Get(fullUrl)
	if err != nil {
		OW.log.Error(fc, ": ", err)
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := err.Error()
		return []byte{}, fmt.Errorf("error: %s", errorMessage)
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			err = json.Unmarshal(body, &errResponse)
			if err == nil {
				OW.log.Error(errResponse)
				return []byte{}, fmt.Errorf("%s", tryAnother)
			}
		}
		OW.log.Error(failedToGetWeather, resp.StatusCode)
		return []byte{}, fmt.Errorf(failedToGetWeather, "%d", resp.StatusCode)
	}
	return body, nil
}

// parseCurrentWeather ...
func (OW *OpenWeatherMapService) parseCurrentWeather(body []byte) (model.WeatherCurrent, error) {
	var weatherData model.WeatherCurrent
	err := json.Unmarshal(body, &weatherData)
	if err != nil {
		OW.log.Error(err, body)
		return model.WeatherCurrent{}, err
	}
	return weatherData, nil
}

// parseForecastWeather ...
func (OW *OpenWeatherMapService) parseForecastWeather(body []byte) (model.WeatherForecast, error) {
	var forecastData model.WeatherForecast
	err := json.Unmarshal(body, &forecastData)
	if err != nil {
		OW.log.Error(err, body)
		return model.WeatherForecast{}, err
	}
	return forecastData, nil
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
