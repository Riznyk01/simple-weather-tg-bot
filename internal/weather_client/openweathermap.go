package weather_client

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/http_client"
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/service/convert"
	"SimpleWeatherTgBot/internal/text"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type OWMService struct {
	httpClient http_client.HTTPClient
	cfg        *config.Config
	log        *logr.Logger
}

func NewOpenWeatherMap(httpClient http_client.HTTPClient, cfg *config.Config, log *logr.Logger) *OWMService {
	return &OWMService{
		httpClient: httpClient,
		cfg:        cfg,
		log:        log,
	}
}

// GetWeatherForecast retrieves the weather forecast based on the provided weather command,
// updates the user's last command, and returns the formatted weather message.
func (OW *OWMService) GetWeatherForecast(us model.UserData) (weatherMessage string, err error) {

	OW.log.V(1).Info("user data: ", "data", us)
	var cityId string
	var weatherData model.WeatherCurrent
	var forecastData model.WeatherForecast

	getWeatherUrl, err := OW.generateWeatherUrl(us)
	if err != nil {
		OW.log.Error(err, text.ErrWhileGeneratingURL)
		return "", err
	}

	r, err := OW.httpClient.Get(getWeatherUrl)
	if err != nil {
		OW.log.Error(err, text.ErrWhileGettingWeather)
		return "", err
	}

	if r.StatusCode != http.StatusOK {
		var errResponse model.ErrorResponse
		if r.StatusCode == http.StatusNotFound {
			decoder := json.NewDecoder(r.Body)
			if err = decoder.Decode(&errResponse); err != nil {
				OW.log.Error(err, text.ErrDecodingJSON)
				return "", err
			} else if err == nil {
				OW.log.V(1).Info(text.ErrFetchingWeather)
				return text.TryAnother, nil
			}
		}
		OW.log.V(1).Info(text.ErrWhileGettingWeather, "StatusCode", r.StatusCode)
		return "", fmt.Errorf(text.FailedToGetWeather)
	}

	decoder := json.NewDecoder(r.Body)
	if us.LastWeatherType == text.CallbackCurrent || us.LastWeatherType == text.CallbackCurrentLocation {
		if err = decoder.Decode(&weatherData); err != nil {
			OW.log.Error(err, text.ErrDecodingJSON)
			return "", err
		}
		weatherMessage, cityId = messageCurrentWeather(weatherData, us.Metric)
	} else if us.LastWeatherType == text.CallbackForecast || us.LastWeatherType == text.CallbackForecastLocation || us.LastWeatherType == text.CallbackTodayLocation || us.LastWeatherType == text.CallbackToday {
		if err = decoder.Decode(&forecastData); err != nil {
			OW.log.Error(err, text.ErrDecodingJSON)
			return "", err
		}
		weatherMessage, cityId = messageForecastWeather(forecastData, us.LastWeatherType, us.Metric)
	}
	return weatherMessage + fmt.Sprintf(text.MoreInfoURLFormat, cityId), nil
}

// generateWeatherUrl ...
func (OW *OWMService) generateWeatherUrl(us model.UserData) (fullWeatherUrl string, err error) {

	OW.log.V(1).Info("command received: ", "command", us)

	weatherUrl := OW.cfg.WeatherApiUrl
	if us.LastWeatherType == text.CallbackCurrent || us.LastWeatherType == text.CallbackCurrentLocation {
		weatherUrl += "weather?"
	} else if us.LastWeatherType == text.CallbackForecast || us.LastWeatherType == text.CallbackForecastLocation || us.LastWeatherType == text.CallbackTodayLocation || us.LastWeatherType == text.CallbackToday {
		weatherUrl += "forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		OW.log.Error(err, text.ErrParsingWeatherURL)
		return "", err
	}

	q := url.Values{}

	if us.LastWeatherType == text.CallbackCurrent || us.LastWeatherType == text.CallbackForecast || us.LastWeatherType == text.CallbackToday {
		q.Add("q", us.City)
	} else if us.LastWeatherType == text.CallbackForecastLocation || us.LastWeatherType == text.CallbackCurrentLocation || us.LastWeatherType == text.CallbackTodayLocation {
		q.Add("lat", us.Lat)
		q.Add("lon", us.Lon)
	}

	q.Add("appid", OW.cfg.WToken)
	if us.Metric {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()

	OW.log.V(1).Info("URL values:", "values", q)
	OW.log.V(1).Info("response URL:", "url", u.String())

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
	return messageWithCurrent + "\n", strconv.Itoa(current.ID)
}

// messageForecastWeather returns a message with weather forecast and city id (in string).
func messageForecastWeather(forecast model.WeatherForecast, command string, metric bool) (message, cityIdStr string) {
	tUnits, wUnits, pUnits := convert.Units(metric)
	// A headers displaying the forecast country and city names,
	// along with units for time, temperature, humidity, pressure, and wind direction.
	headerPlace := fmt.Sprintf("<b>%s %s\n\n</b>", forecast.City.Country, forecast.City.Name)
	headerUnits := fmt.Sprintf("[TIME] [%s] [%s] [%s] [%s, dir.]\n", tUnits, "💧", pUnits, wUnits)
	message += headerPlace
	// Iterate through forecast entries to construct messages for each time period.
	//for i, e := range forecast.List {
	for i := 0; i < len(forecast.List); i++ {
		forecastTime := time.Unix(int64(forecast.List[i].Dt), 0).
			In(time.FixedZone("Custom Timezone", forecast.City.Timezone))
		// Check if the user needs today's forecast.
		// Breaks if current entry is the next day.
		if (command == text.CallbackTodayLocation || command == text.CallbackToday) && i != 0 {
			lastEntrysDay, _ := strconv.Atoi(time.Unix(int64(forecast.List[i-1].Dt), 0).
				In(time.FixedZone("Custom Timezone", forecast.City.Timezone)).Format("02"))
			currentDay, _ := strconv.Atoi(forecastTime.Format("02"))
			if currentDay > lastEntrysDay {
				break
			}
		}
		hours, day := forecastTime.Format("15"), forecastTime.Format("02")
		windSpeedForecast := convert.WindSpeed(forecast.List[i].Wind.Speed, metric)
		// Display the date for each day in the format:
		// 🗓 31 January (Wednesday) along with the header containing units.
		entrysHour, _ := strconv.Atoi(hours)
		if entrysHour < 3 || i == 0 {
			message += fmt.Sprintf("<b>🗓 %s %s (%s)</b>\n",
				day, forecastTime.Month().String(), forecastTime.Weekday().String()) + headerUnits
		}
		// Hourly forecast.
		message += fmt.Sprintf("%s:00 %s %+d %d%% %d [%.1f %s]\n",
			hours,
			convert.AddIcon(forecast.List[i].Weather[0].Description, false),
			convert.KelvinToFahrenheitAndRound(forecast.List[i].Main.Temp, metric),
			forecast.List[i].Main.Humidity,
			convert.Pressure(forecast.List[i].Main.Pressure, metric),
			windSpeedForecast, convert.DegsToDirIcon(forecast.List[i].Wind.Deg))
		// Insert a newline between days at the end of each day.
		if entrysHour == 23 && i != len(forecast.List) {
			message += "\n"
		}
	}
	return message, strconv.Itoa(forecast.City.ID)
}
