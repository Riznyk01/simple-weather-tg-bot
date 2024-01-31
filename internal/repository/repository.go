package repository

import (
	"SimpleWeatherTgBot/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	SetUserMeasurementSystem(id int64, system bool) error
	SetUserLastInputCity(id int64, city string) error
	SetUserLastInputLocation(id int64, lat, lon string) error
	SetUserLastWeatherCommand(id int64, last string) error
	GetUserById(id int64) (model.UserData, error)
	CreateUserById(userId int64) error
}

type Repository struct {
	UserRepository
}

func NewRepository(log *logrus.Logger, db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(log, db),
	}
}
