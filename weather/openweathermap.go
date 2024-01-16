package weather

import "SimpleWeatherTgBot/repository"

type OpenWeatherMapService struct {
	repo *repository.Repository
}

func NewOpenWeatherMap(repo *repository.Repository) *OpenWeatherMapService {
	return &OpenWeatherMapService{
		repo: repo,
	}
}

func (OW *OpenWeatherMapService) SetSystem(id int64, system bool) {
	OW.repo.SetSystem(id, system)
}
func (OW *OpenWeatherMapService) SetCity(id int64, city string) {

}
func (OW *OpenWeatherMapService) SetLocation(id int64, lat, lon string) {

}
func (OW *OpenWeatherMapService) SetLast(id int64, last string) {

}
func (OW *OpenWeatherMapService) GetSystem(id int64) (system bool) {
	return true
}
func (OW *OpenWeatherMapService) GetCity(id int64) string {
	return ""
}
func (OW *OpenWeatherMapService) GetLat(id int64) string {
	return ""
}
func (OW *OpenWeatherMapService) GetLon(id int64) string {
	return ""
}
func (OW *OpenWeatherMapService) GetLast(id int64) string {
	return ""
}
func (OW *OpenWeatherMapService) Exists(id int64) bool {
	return true
}

//
//var weatherData types.WeatherResponse
//var forecastData types.WeatherResponse5d3h
//var temperatureUnits, windUnits, pressureUnits string
//
//const (
//	apiWeatherUrl = "https://api.openweathermap.org/data/2.5/"
//)
//
////TODO: add logger
//
//// Returns a complete weather message.
//func (wc *WClient) LastWeather(chatId int64) (string, error) {
//
//	return "", nil
//}
//
//func GetWeather(fullUrlGet, forecastType string, metric bool) (weatherMessage string, err error) {
//	var cityIdString string
//
//	resp, err := http.Get(fullUrlGet)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		errorMessage := err.Error()
//		return "", fmt.Errorf("error: %s", errorMessage)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		if resp.StatusCode == http.StatusNotFound {
//			var errorResponse struct {
//				Cod     string `json:"cod"`
//				Message string `json:"message"`
//			}
//			err = json.Unmarshal(body, &errorResponse)
//			if err == nil {
//				return "", fmt.Errorf("%s. Try another city name.", errorResponse.Message)
//			}
//		}
//		return "", fmt.Errorf("Failed to get weather data. Status code: %d", resp.StatusCode)
//	}
//	//Getting units of measurement
//	temperatureUnits, windUnits, pressureUnits = units(metric)
//	if forecastType == types.CommandCurrent || forecastType == types.CommandCurrentLocation {
//		err = json.Unmarshal(body, &weatherData)
//		if err != nil {
//			return "", err
//		}
//		weatherMessage, cityIdString = messageCurrentWeather(weatherData, metric)
//	} else if forecastType == types.CommandForecast || forecastType == types.CommandForecastLocation {
//		err = json.Unmarshal(body, &forecastData)
//		if err != nil {
//			return "", err
//		}
//		weatherMessage, cityIdString = messageForecastWeather(forecastData, metric)
//	}
//	more := fmt.Sprintf("\n\n<a href=\"https://openweathermap.org/city/%s\">🌐 More</a>", cityIdString)
//	return weatherMessage + more, nil
//}
//
//func GenerateWeatherUrl(weatherParam map[string]string, tWeather, forecastType string, metricUnits bool) (string, error) {
//	var weatherUrl string
//
//	if forecastType == types.CommandCurrent || forecastType == types.CommandCurrentLocation {
//		weatherUrl = apiWeatherUrl + "weather?"
//	} else if forecastType == types.CommandForecast || forecastType == types.CommandForecastLocation {
//		weatherUrl = apiWeatherUrl + "forecast?"
//	}
//
//	u, err := url.Parse(weatherUrl)
//	if err != nil {
//		return "", err
//	}
//
//	q := url.Values{}
//	if _, ex := weatherParam["city"]; ex {
//		q.Add("q", weatherParam["city"])
//	} else if _, exLat := weatherParam["lat"]; exLat {
//		q.Add("lat", weatherParam["lat"])
//		q.Add("lon", weatherParam["lon"])
//	}
//	q.Add("appid", tWeather)
//	if metricUnits {
//		q.Add("units", "metric")
//	}
//	u.RawQuery = q.Encode()
//	fullUrlGet := u.String()
//	return fullUrlGet, nil
//}
//
//// Returns units based on the metric system.
//func units(metricUnits bool) (tempUnits, windUnits, pressureUnits string) {
//	if metricUnits {
//		tempUnits = " °C"
//		windUnits = " m/s"
//		pressureUnits = " mmHg"
//	} else {
//		tempUnits = " °F"
//		windUnits = " mph"
//		pressureUnits = " inHg"
//	}
//	return tempUnits, windUnits, pressureUnits
//}
//
//// Returns a message with current weather and city id (in string).
//func messageCurrentWeather(weatherData types.WeatherResponse, metric bool) (userMessageCurrent, cityIdStr string) {
//	pressure := utils.PressureConverting(float64(weatherData.Main.Pressure), metric)
//	windSpeed := weatherData.Wind.Speed
//	//Converting to miles per hour if non-metric
//	if !metric {
//		windSpeed = utils.ToMilesPerHour(weatherData.Wind.Speed)
//		pressure = utils.PressureConverting(float64(weatherData.Main.Pressure), metric)
//	}
//	userMessageCurrent = fmt.Sprintf("<b>%s %s</b> %s\n\n 🌡 %+d%s (Feel %+d%s) 💧 %d%%  \n\n 📉 %+d%s ️ 📈 %+d%s \n%.0f %s %.2f%s %s \n\n🌅  %s 🌉  %s",
//		weatherData.Sys.Country,
//		weatherData.Name,
//		utils.ReplaceWeatherPlusIcons(weatherData.Weather[0].Description),
//		utils.TemperatureConverting(weatherData.Main.Temp, metric),
//		temperatureUnits,
//		utils.TemperatureConverting(weatherData.Main.FeelsLike, metric),
//		temperatureUnits,
//		weatherData.Main.Humidity,
//		utils.TemperatureConverting(weatherData.Main.TempMin, metric),
//		temperatureUnits,
//		utils.TemperatureConverting(weatherData.Main.TempMax, metric),
//		temperatureUnits,
//		pressure,
//		pressureUnits,
//		windSpeed,
//		windUnits,
//		utils.DegreesToDirectionIcon(weatherData.Wind.Deg),
//		utils.TimeStampToHuman(weatherData.Sys.Sunrise, weatherData.Timezone, "15:04"),
//		utils.TimeStampToHuman(weatherData.Sys.Sunset, weatherData.Timezone, "15:04"))
//
//	return userMessageCurrent, strconv.Itoa(weatherData.ID)
//}
//
//// Returns a message with weather forecast and city id (in string).
//func messageForecastWeather(forecastData types.WeatherResponse5d3h, metric bool) (userMessageForecast, cityIdStr string) {
//	// Creating a string to display the country and city names
//	userMessageForecast = fmt.Sprintf("<b>%s %s\n\n</b>", forecastData.City.Country, forecastData.City.Name)
//	// Constructing the date display, including day, month, and day of the week,
//	// to be inserted into the user message about the weather.
//	userMessageForecast += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", utils.TimeStampToHuman(forecastData.List[0].Dt, forecastData.City.Timezone, "02"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "m"), utils.TimeStampToInfo(forecastData.List[0].Dt, forecastData.City.Timezone, "d"))
//	messageHeader := fmt.Sprintf("[%s] [---] [%s] [%s] [%s] [%s]\n",
//		"h:m",
//		temperatureUnits,
//		"%",
//		pressureUnits,
//		windUnits,
//	)
//
//	userMessageForecast += messageHeader
//
//	for ind, entry := range forecastData.List {
//		hours := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "15")
//		dayNum := utils.TimeStampToHuman(entry.Dt, forecastData.City.Timezone, "02")
//		dayOfWeek := utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "d")
//		if hours == "01" || hours == "02" && ind > 0 {
//			// Constructing the date display, including day, month, and day of the week,
//			// to be inserted into the user message about the weather.
//			userMessageForecast += fmt.Sprintf("<b>🗓%s %s (%s)</b>\n", dayNum, utils.TimeStampToInfo(entry.Dt, forecastData.City.Timezone, "m"), dayOfWeek)
//			userMessageForecast += messageHeader
//		}
//
//		windSpeedForecast := entry.Wind.Speed
//		// Converting to miles per hour if non-metric
//		if !metric {
//			windSpeedForecast = utils.ToMilesPerHour(entry.Wind.Speed)
//		}
//
//		userMessageForecast += fmt.Sprintf("%s %s %+d %d%% %.1f %.1f %s\n",
//			hours+":00",
//			utils.ReplaceWeatherToIcons(entry.Weather[0].Description),
//			utils.TemperatureConverting(entry.Main.Temp, metric),
//			entry.Main.Humidity,
//			utils.PressureConverting(float64(entry.Main.Pressure), metric),
//			windSpeedForecast,
//			utils.DegreesToDirectionIcon(entry.Wind.Deg),
//		)
//
//		if hours == "21" || hours == "22" || hours == "23" {
//			userMessageForecast += "\n"
//		}
//
//	}
//	return userMessageForecast, strconv.Itoa(forecastData.City.ID)
//}
