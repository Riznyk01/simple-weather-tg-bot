package repository

import "SimpleWeatherTgBot/types"

type MemoryStorage struct {
	data map[int64]types.UserData
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[int64]types.UserData),
	}
}
