package admin

import (
	"SimpleWeatherTgBot/config"
	"SimpleWeatherTgBot/repository"
	"github.com/sirupsen/logrus"
)

type AdminControlMemoryStorage struct {
	repo *repository.Repository
	log  *logrus.Logger
	cfg  *config.Config
}

func NewAdminControlMemoryStorage(repo *repository.Repository, log *logrus.Logger, cfg *config.Config) *AdminControlMemoryStorage {
	return &AdminControlMemoryStorage{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

func (uc *AdminControlMemoryStorage) BanUser(userId int64) error {
	return uc.repo.BanUser(userId)
}

func (uc *AdminControlMemoryStorage) UnbanUser(userId int64) error {
	return uc.repo.UnbanUser(userId)
}

func (uc *AdminControlMemoryStorage) CheckPass(pass string, passCfg string) bool {
	return pass == passCfg
}
