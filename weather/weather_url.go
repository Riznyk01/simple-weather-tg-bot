package weather

import (
	"SimpleWeatherTgBot/lib/e"
	"net/url"
)

const (
	apiWeatherUrl = "https://api.openweathermap.org/data/2.5/"
)

func WeatherUrlByCity(city, tWeather, forecastType string, metricUnits bool) (string, bool, error) {
	var weatherUrl string
	if forecastType == "current" {
		weatherUrl = apiWeatherUrl + "weather?"
	} else if forecastType == "5-days forecast" {
		weatherUrl = apiWeatherUrl + "forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		return "", false, e.Wrap("", err)
	}

	q := url.Values{}
	q.Add("q", city)
	q.Add("appid", tWeather)
	if metricUnits {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	return fullUrlGet, metricUnits, nil
}

func WeatherUrlByLocation(latStr, lonStr, tWeather, forecastType string, metricUnits bool) (string, bool, error) {
	var weatherUrl string

	if forecastType == "current 📍" {
		weatherUrl = apiWeatherUrl + "weather?"
	} else if forecastType == "5-days forecast 📍" {
		weatherUrl = apiWeatherUrl + "forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		return "", false, e.Wrap("", err)
	}

	q := url.Values{}
	q.Add("lat", latStr)
	q.Add("lon", lonStr)
	q.Add("appid", tWeather)
	if metricUnits {
		q.Add("units", "metric")
	}
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	return fullUrlGet, metricUnits, nil
}
