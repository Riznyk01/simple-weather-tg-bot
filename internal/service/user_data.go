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
	fc := "SetCity"
	err = UP.repo.SetCity(chatId, city)
	if err != nil {
		UP.log.Errorf("%s: %v", fc, err)
		return err
	}
	return nil
}

// SetLocation ...
func (UP *UserDataService) SetLocation(chatId int64, lat, lon string) error {
	fc := "SetLocation"
	err := UP.repo.SetLocation(chatId, lat, lon)
	if err != nil {
		UP.log.Errorf("%s: %v", fc, err)
		return err
	}
	return nil
}

func (UP *UserDataService) SetLastWeatherCommand(chatId int64, last string) error {
	fc := "SetLastWeatherCommand"
	err := UP.repo.SetLastWeatherCommand(chatId, last)
	if err != nil {
		UP.log.Errorf("%s: %v", fc, err)
		return err
	}
	return nil
}

// GetUser ...
func (UP *UserDataService) GetUser(userId int64) (model.UserData, error) {
	fc := "GetUser"
	user, err := UP.repo.GetUser(userId)
	if err != nil {
		UP.log.Errorf("%s: %v", fc, err)
		return model.UserData{}, nil
	}
	return user, nil
}
