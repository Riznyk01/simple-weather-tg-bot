package weather

import (
	"SimpleWeatherTgBot/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetWeather(city, tWeather string) (types.WeatherResponse, error) {

	weatherUrl := "https://api.openweathermap.org/data/2.5/weather?"

	u, err := url.Parse(weatherUrl)
	if err != nil {
		fmt.Println("Error parsing URL (getWeather):", err)
		return types.WeatherResponse{}, err
	}
	q := url.Values{}
	q.Add("q", city)
	q.Add("appid", tWeather)
	q.Add("units", "metric")
	u.RawQuery = q.Encode()
	fullUrlGet := u.String()
	fmt.Println(fullUrlGet)
	resp, err := http.Get(fullUrlGet)
	if err != nil {
		return types.WeatherResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("resp error:", err)
		return types.WeatherResponse{}, err
	}
	var weatherResponse types.WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		fmt.Println("getWeather func err:", err)
		return types.WeatherResponse{}, err
	}
	return weatherResponse, nil
}
