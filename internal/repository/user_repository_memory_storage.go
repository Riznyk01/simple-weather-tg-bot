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

// SetUserMeasurementSystem sets user's system of measurement.
func (u *Memory) SetUserMeasurementSystem(id int64, system bool) error {
	currentData := u.memory.data[id]
	currentData.Metric = system
	u.memory.data[id] = currentData
	fmt.Println(u.memory.data)
	return nil
}

// SetUserLastInputCity sets the user's last input city for weather forecast.
func (u *Memory) SetUserLastInputCity(id int64, city string) error {
	currentData := u.memory.data[id]
	currentData.City = city
	u.memory.data[id] = currentData
	return nil
}

// SetUserLastInputLocation sets the user's last input location for weather forecast.
func (u *Memory) SetUserLastInputLocation(id int64, lat, lon string) error {
	currentData := u.memory.data[id]
	currentData.Lat = lat
	currentData.Lon = lon
	u.memory.data[id] = currentData
	return nil
}

// SetUserLastWeatherCommand sets the user's last input forecast type.
func (u *Memory) SetUserLastWeatherCommand(id int64, last string) error {
	currentData := u.memory.data[id]
	currentData.Last = last
	u.memory.data[id] = currentData
	return nil
}

// GetUserById gets the user's data from the database.
func (u *Memory) GetUserById(id int64) (model.UserData, error) {
	userData, _ := u.memory.data[id]
	return userData, nil
}

// CreateUserById creates a user in the database.
func (u *Memory) CreateUserById(userId int64) error {
	//
	return nil
}
