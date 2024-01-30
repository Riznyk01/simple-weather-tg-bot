package repository

import (
	"SimpleWeatherTgBot/internal/model"
)

type UserRepository interface {
	SetSystem(id int64, system bool) error
	SetCity(id int64, city string) error
	SetLocation(id int64, lat, lon string) error
	SetLastWeatherCommand(id int64, last string) error
	GetUser(id int64) (model.UserData, error)
}

type Repository struct {
	UserRepository
}

func NewRepository(memoryStor *MemoryStorage) *Repository {
	return &Repository{
		UserRepository: NewUserMemoryStorage(memoryStor),
	}
}
