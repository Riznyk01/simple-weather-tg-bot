package service

import (
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/repository"
	"github.com/sirupsen/logrus"
)

type UserDataService struct {
	log  *logrus.Logger
	repo *repository.Repository
}

func NewUserPreferencesService(log *logrus.Logger, repo *repository.Repository) *UserDataService {
	return &UserDataService{
		log:  log,
		repo: repo,
	}
}

// SetSystem ...
func (UP *UserDataService) SetSystem(chatId int64, command string) (err error) {
	fc := "SetSystem"
	m := false
	if command == model.CommandMetricUnits {
		m = true
	}
	err = UP.repo.SetSystem(chatId, m)
	if err != nil {
		UP.log.Errorf("%s: %v", fc, err)
		return err
	}
	return nil
}

// SetCity ...
func (UP *UserDataService) SetCity(chatId int64, city string) (err error) {
	return UP.repo.SetCity(chatId, city)
}

// SetLocation ...
func (UP *UserDataService) SetLocation(chatId int64, lat, lon string) error {
	return UP.repo.SetLocation(chatId, lat, lon)
}

func (UP *UserDataService) SetLastWeatherCommand(chatId int64, last string) error {
	return UP.repo.SetLastWeatherCommand(chatId, last)
}

// GetUserById ...
func (UP *UserDataService) GetUserById(userId int64) (model.UserData, error) {
	return UP.repo.GetUserById(userId)
}

func (UP *UserDataService) CreateUser(userId int64) error {
	return UP.repo.CreateUser(userId)
}
