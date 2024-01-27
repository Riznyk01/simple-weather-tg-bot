package weather_service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/repository"
	"github.com/sirupsen/logrus"
)

type WeatherControl interface {
	GetWeatherForecast(id int64, command string) (string, error)
}

type WeatherService struct {
	WeatherControl
}

func NewWeatherService(repo *repository.Repository, cfg *config.Config, log *logrus.Logger) *WeatherService {
	return &WeatherService{
		WeatherControl: NewOpenWeatherMap(
			repo,
			cfg,
			log),
	}
}
