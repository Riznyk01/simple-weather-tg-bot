package service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/http_client"
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/repository"
	"github.com/sirupsen/logrus"
)

type UserData interface {
	SetUserMeasurementSystem(id int64, command string) error
	SetUserLastInputCity(id int64, city string) error
	SetUserLastInputLocation(id int64, lat, lon string) error
	SetUserLastWeatherCommand(id int64, last string) error
	GetUserById(id int64) (model.UserData, error)
	CreateUserById(userId int64) error
}

type WeatherApi interface {
	GetWeatherForecast(user model.UserData) (string, error)
}

type Service struct {
	UserData
	WeatherApi
}

func NewService(repo *repository.Repository, cfg *config.Config, log *logrus.Logger, httpClient http_client.HTTPClient) *Service {
	return &Service{
		UserData:   NewUserPreferencesService(log, repo),
		WeatherApi: NewOpenWeatherMap(httpClient, cfg, log),
	}
}
