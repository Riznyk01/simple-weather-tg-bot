package weather_client

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/http_client"
	"SimpleWeatherTgBot/internal/model"
	"github.com/go-logr/logr"
)

type Client interface {
	GetWeatherForecast(user model.UserData) (string, error)
}

type Weather struct {
	Client
}

func NewWeather(httpClient http_client.HTTPClient, cfg *config.Config, log *logr.Logger) *Weather {
	return &Weather{
		Client: NewOpenWeatherMap(httpClient, cfg, log),
	}
}
