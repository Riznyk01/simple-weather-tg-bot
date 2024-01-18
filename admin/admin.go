package admin

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/repository"
	"github.com/sirupsen/logrus"
)

type AdminControl interface {
	BanUser(id int64) error
	UnbanUser(id int64) error
}

type AdminService struct {
	AdminControl
}

func NewAdminService(repo *repository.Repository, log *logrus.Logger, cfg *config.Config) *AdminService {
	return &AdminService{
		AdminControl: NewAdminControlMemoryStorage(repo, log, cfg),
	}
}
