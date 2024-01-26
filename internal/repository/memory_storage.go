package repository

import (
	"SimpleWeatherTgBot/internal/model"
)

type MemoryStorage struct {
	data map[int64]model.UserData
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[int64]model.UserData),
	}
}
