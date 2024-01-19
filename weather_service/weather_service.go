package weather_service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/repository"
	"github.com/sirupsen/logrus"
)

type WeatherControl interface {
	SetSystem(id int64, system bool) error
	SetCity(id int64, city string) error
	SetLocation(id int64, lat, lon string) error
	SetLast(id int64, command string) (string, error)
	GetSystem(id int64) (bool, error)
	GetCity(id int64) (string, error)
	GetLocation(id int64) (string, string, error)
	GetLast(id int64) (string, error)
	AddRequestsCount(id int64) (int, error)
}

type WeatherService struct {
	WeatherControl
}

func NewWClient(repo *repository.Repository, cfg *config.Config, log *logrus.Logger) *WeatherService {
	return &WeatherService{
		WeatherControl: NewOpenWeatherMap(
			repo,
			cfg,
			log),
	}
}
