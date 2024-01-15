package weather

import (
	"SimpleWeatherTgBot/types"
	"net/url"
)

const (
	apiWeatherUrl = "https://api.openweathermap.org/data/2.5/"
)

// TODO: add logger
func GenerateWeatherUrl(weatherParam map[string]string, tWeather, forecastType string, metricUnits bool) (string, error) {
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
	if _, ex := weatherParam["city"]; ex {
		q.Add("q", weatherParam["city"])
	} else if _, exLat := weatherParam["lat"]; exLat {
		q.Add("lat", weatherParam["lat"])
		q.Add("lon", weatherParam["lon"])
	}
	q.Add("appid", tWeather)
	if metricUnits {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	return fullUrlGet, nil
}
