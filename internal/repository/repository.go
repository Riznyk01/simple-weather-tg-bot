package repository

import "errors"

var ErrItemIsEmpty = errors.New("item is empty")

type WeatherControl interface {
	SetSystem(id int64, system bool) error
	SetCity(id int64, city string) error
	SetLocation(id int64, lat, lon string) error
	SetLast(id int64, last string) error
	GetSystem(id int64) (bool, error)
	GetCity(id int64) (string, error)
	GetLocation(id int64) (string, string, error)
	GetLast(id int64) (string, error)
}

type Repository struct {
	WeatherControl
}

func NewRepository(memoryStor *MemoryStorage) *Repository {
	return &Repository{
		WeatherControl: NewWeatherControlMemoryStorage(memoryStor),
	}
}
