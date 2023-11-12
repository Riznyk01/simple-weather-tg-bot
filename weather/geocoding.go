package weather

import (
	"SimpleWeatherTgBot/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CoordinatesByLocationName(directGeoUrl, endOfDirectGeoUrl string, update types.Update) ([]types.Geocoding, error) {
	city := update.Message.Text
	resp, err := http.Get(directGeoUrl + city + endOfDirectGeoUrl)
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
