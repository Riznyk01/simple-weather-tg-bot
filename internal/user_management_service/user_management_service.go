package user_management_service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/repository"
	"github.com/sirupsen/logrus"
)

type UserControl interface {
	SetLocation(id int64, lat, lon string) error
	GetLocation(id int64) (string, string, error)
	SetLastWeatherCommand(id int64, command string) error
	GetLastWeatherCommand(id int64) (string, error)
	SetSystem(id int64, command string) error
	GetSystem(id int64) (bool, error)
	SetCity(id int64, city string) error
	GetCity(id int64) (string, error)
}

type UserService struct {
	UserControl
}

func NewUserService(repo *repository.Repository, cfg *config.Config, log *logrus.Logger) *UserService {
	return &UserService{
		UserControl: NewUserControlImpl(
			repo,
			cfg,
			log),
	}
}
