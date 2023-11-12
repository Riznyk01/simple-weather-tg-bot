package weather

import (
	"SimpleWeatherTgBot/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func CoordinatesByLocationName(tWeather string, update types.Update) ([]types.Geocoding, error) {
	city := update.Message.Text

	directGeoUrl := "http://api.openweathermap.org/geo/1.0/direct?q="
	limit := "1"

	u, err := url.Parse(directGeoUrl)
	if err != nil {
		fmt.Println("Error parsing URL (getWeather):", err)
		return []types.Geocoding{}, err
	}
	q := url.Values{}
	q.Add("q", city)
	q.Add("limit", limit)
	q.Add("appid", tWeather)

	u.RawQuery = q.Encode()
	fullUrlLocation := u.String()

	resp, err := http.Get(fullUrlLocation)
	if err != nil {
		fmt.Println(err)
		return []types.Geocoding{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return []types.Geocoding{}, err
	}
	var geocoding []types.Geocoding

	err = json.Unmarshal(body, &geocoding)
	if err != nil {
		fmt.Println("CoordinatesByLocationName func err:", err)
		return []types.Geocoding{}, err
	}
	return geocoding, nil
}
