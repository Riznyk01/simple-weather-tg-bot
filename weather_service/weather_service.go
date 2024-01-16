package weather_service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/repository"
	"github.com/sirupsen/logrus"
)

type WeatherUserControl interface {
	SetSystem(id int64, system bool)
	SetCity(id int64, city string)
	SetLocation(id int64, lat, lon string)
	SetLast(id int64, last string) (string, error)
	GetSystem(id int64) (bool, error)
	GetCity(id int64) string
	GetLat(id int64) string
	GetLon(id int64) string
	GetLast(id int64) (weatherMessage string, err error)
	Exists(id int64) bool
}

type WeatherService struct {
	WeatherUserControl
}

func NewWClient(repo *repository.Repository, cfg *config.Config, log *logrus.Logger) *WeatherService {
	return &WeatherService{
		WeatherUserControl: NewOpenWeatherMap(
			repo,
			cfg,
			log),
	}
}
