package weather

import (
	"fmt"
	"log"
	"net/url"
)

func WeatherNowUrlByCity(city, tWeather string) string {
	weatherUrl := "https://api.openweathermap.org/data/2.5/weather?"

	u, err := url.Parse(weatherUrl)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		fmt.Errorf("error: %s", errorMessage)
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
func Weather5d3hUrlByCity(city, tWeather string) string {
	weatherUrl := "https://api.openweathermap.org/data/2.5/forecast?"

	u, err := url.Parse(weatherUrl)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		fmt.Errorf("error: %s", errorMessage)
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
func WeatherNowUrlByLocation(latStr, lonStr, tWeather string) string {
	weatherUrl := "https://api.openweathermap.org/data/2.5/weather?"

	u, err := url.Parse(weatherUrl)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		fmt.Errorf("error: %s", errorMessage)
		return ""
	}
	q := url.Values{}
	q.Add("lat", latStr)
	q.Add("lon", lonStr)
	q.Add("appid", tWeather)
	q.Add("units", "metric")
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	return fullUrlGet
}

func Weather5d3hUrlByLocation(latStr, lonStr, tWeather string) string {
	weatherUrl := "https://api.openweathermap.org/data/2.5/forecast?"

	u, err := url.Parse(weatherUrl)
	if err != nil {
		errorMessage := err.Error()
		log.Println("Error: ", errorMessage)
		fmt.Errorf("error: %s", errorMessage)
		return ""
	}
	q := url.Values{}
	q.Add("lat", latStr)
	q.Add("lon", lonStr)
	q.Add("appid", tWeather)
	q.Add("units", "metric")
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	return fullUrlGet
}
