package repository

import (
	"SimpleWeatherTgBot/internal/model"
	"fmt"
)

type MemoryStorage struct {
	data map[int64]model.UserData
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[int64]model.UserData),
	}
}

type Memory struct {
	memory *MemoryStorage
}

func NewUserMemoryStorage(memory *MemoryStorage) *Memory {
	return &Memory{
		memory: memory,
	}
}

// SetSystem sets user's system of measurement.
func (u *Memory) SetSystem(id int64, system bool) error {
	currentData := u.memory.data[id]
	currentData.Metric = system
	u.memory.data[id] = currentData
	fmt.Println(u.memory.data)
	return nil
}

// SetCity gets user's last city.
func (u *Memory) SetCity(id int64, city string) error {
	currentData := u.memory.data[id]
	currentData.City = city
	u.memory.data[id] = currentData
	return nil
}

// SetLocation gets user's last location.
func (u *Memory) SetLocation(id int64, lat, lon string) error {
	currentData := u.memory.data[id]
	currentData.Lat = lat
	currentData.Lon = lon
	u.memory.data[id] = currentData
	return nil
}

// SetLastWeatherCommand sets last users forecast type.
func (u *Memory) SetLastWeatherCommand(id int64, last string) error {
	currentData := u.memory.data[id]
	currentData.Last = last
	u.memory.data[id] = currentData
	return nil
}

// GetUser ...
func (u *Memory) GetUser(id int64) (model.UserData, error) {
	userData, _ := u.memory.data[id]
	return userData, nil
}
