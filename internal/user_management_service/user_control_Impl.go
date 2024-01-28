package user_management_service

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/repository"
	"github.com/sirupsen/logrus"
)

type UserControlService struct {
	repo *repository.Repository
	cfg  *config.Config
	log  *logrus.Logger
}

func NewUserControlImpl(repo *repository.Repository, cfg *config.Config, log *logrus.Logger) *UserControlService {
	return &UserControlService{
		repo: repo,
		cfg:  cfg,
		log:  log,
	}
}

// SetSystem ...
func (UC *UserControlService) SetSystem(chatId int64, command string) error {
	if command == model.CommandMetricUnits {
		return UC.repo.SetSystem(chatId, true)
	}
	return UC.repo.SetSystem(chatId, false)
}

// SetCity ...
func (UC *UserControlService) SetCity(chatId int64, city string) error {
	return UC.repo.SetCity(chatId, city)
}

// SetLocation ...
func (UC *UserControlService) SetLocation(chatId int64, lat, lon string) error {
	fc := "SetLocation"
	err := UC.repo.SetLocation(chatId, lat, lon)
	if err != nil {
		UC.log.Infof("%s: %v", fc, err)
	}
	return nil
}

// GetSystem ...
func (UC *UserControlService) GetSystem(chatId int64) (bool, error) {
	return UC.repo.GetSystem(chatId)
}

// GetCity ...
func (UC *UserControlService) GetCity(chatId int64) (string, error) {
	return UC.repo.GetCity(chatId)
}

// GetLocation ...
func (UC *UserControlService) GetLocation(chatId int64) (string, string, error) {
	return UC.repo.GetLocation(chatId)
}
