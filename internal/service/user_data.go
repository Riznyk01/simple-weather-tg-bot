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

// SetUserMeasurementSystem sets user's system of measurement.
func (UP *UserDataService) SetUserMeasurementSystem(chatId int64, command string) (err error) {
	fc := "SetUserMeasurementSystem"
	m := false
	if command == model.CommandMetricUnits {
		m = true
	}
	err = UP.repo.SetUserMeasurementSystem(chatId, m)
	if err != nil {
		UP.log.Errorf("%s: %v", fc, err)
		return err
	}
	return nil
}

// SetUserLastInputCity sets the user's last input city for weather forecast.
func (UP *UserDataService) SetUserLastInputCity(chatId int64, city string) (err error) {
	return UP.repo.SetUserLastInputCity(chatId, city)
}

// SetUserLastInputLocation sets the user's last input location for weather forecast.
func (UP *UserDataService) SetUserLastInputLocation(chatId int64, lat, lon string) error {
	return UP.repo.SetUserLastInputLocation(chatId, lat, lon)
}

// SetUserLastWeatherCommand sets the user's last input forecast type.
func (UP *UserDataService) SetUserLastWeatherCommand(chatId int64, last string) error {
	return UP.repo.SetUserLastWeatherCommand(chatId, last)
}

// GetUserById gets the user's data from the database.
func (UP *UserDataService) GetUserById(userId int64) (model.UserData, error) {
	return UP.repo.GetUserById(userId)
}

// CreateUserById creates a user in the database.
func (UP *UserDataService) CreateUserById(userId int64) error {
	return UP.repo.CreateUserById(userId)
}
