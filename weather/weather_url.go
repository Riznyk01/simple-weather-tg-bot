package weather

import (
	"SimpleWeatherTgBot/lib/e"
	"log"
	"net/url"
)

func WeatherUrlByCity(city, tWeather, forecastType string) string {
	var weatherUrl string

	if forecastType == "current" {
		weatherUrl = "https://api.openweathermap.org/data/2.5/weather?"
	} else if forecastType == "5-days forecast" {
		weatherUrl = "https://api.openweathermap.org/data/2.5/forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		log.Println(e.Wrap("", err))
		return ""
	}

	q := url.Values{}
	q.Add("q", city)
	q.Add("appid", tWeather)
	q.Add("units", "metric")
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	return fullUrlGet
}

func WeatherUrlByLocation(latStr, lonStr, tWeather, forecastType string) (string, error) {
	var weatherUrl string

	if forecastType == "current 📍" {
		weatherUrl = "https://api.openweathermap.org/data/2.5/weather?"
	} else if forecastType == "5-days forecast 📍" {
		weatherUrl = "https://api.openweathermap.org/data/2.5/forecast?"
	}

	u, err := url.Parse(weatherUrl)
	if err != nil {
		return "", e.Wrap("", err)
	}

	q := url.Values{}
	q.Add("lat", latStr)
	q.Add("lon", lonStr)
	q.Add("appid", tWeather)
	q.Add("units", "metric")
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	return fullUrlGet, nil
}
