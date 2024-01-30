package repository

import (
	"SimpleWeatherTgBot/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	SetSystem(id int64, system bool) error
	SetCity(id int64, city string) error
	SetLocation(id int64, lat, lon string) error
	SetLastWeatherCommand(id int64, last string) error
	GetUserById(id int64) (model.UserData, error)
	CreateUser(userId int64) error
}

type Repository struct {
	UserRepository
}

func NewRepository(log *logrus.Logger, db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(log, db),
	}
}
