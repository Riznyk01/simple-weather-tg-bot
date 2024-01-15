package weather

import (
	"SimpleWeatherTgBot/types"
	"net/url"
)

const (
	apiWeatherUrl = "https://api.openweathermap.org/data/2.5/"
)

// TODO: add logger
func GenerateWeatherUrl(weatherParam []string, tWeather, forecastType string, metricUnits bool) (string, error) {
	var weatherUrl string

	if forecastType == types.CommandCurrent || forecastType == types.CommandCurrentLocation {
		weatherUrl = apiWeatherUrl + "weather?"
	} else if forecastType == types.CommandForecast || forecastType == types.CommandForecastLocation {
		weatherUrl = apiWeatherUrl + "forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	if len(weatherParam) == 1 {
		q.Add("q", weatherParam[0])
	} else if len(weatherParam) == 2 {
		q.Add("lat", weatherParam[0])
		q.Add("lon", weatherParam[1])
	}
	q.Add("appid", tWeather)
	if metricUnits {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	return fullUrlGet, nil
}
