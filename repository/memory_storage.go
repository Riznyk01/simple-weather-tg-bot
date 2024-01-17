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

func (u *MemoryStorage) SetSystem(id int64, system bool) error {
	currentData := u.data[id]
	currentData.Metric = system
	u.data[id] = currentData
	return nil
}

func (u *MemoryStorage) SetCity(id int64, city string) error {
	currentData := u.data[id]
	currentData.City = city
	u.data[id] = currentData
	return nil
}

func (u *MemoryStorage) SetLocation(id int64, lat, lon string) error {
	currentData := u.data[id]
	currentData.Lat = lat
	currentData.Lon = lon
	u.data[id] = currentData
	return nil
}

// Set last users forecast type
func (u *MemoryStorage) SetLast(id int64, last string) error {
	currentData := u.data[id]
	currentData.Last = last
	u.data[id] = currentData
	return nil
}

func (u *MemoryStorage) GetSystem(id int64) (bool, error) {
	return u.data[id].Metric, nil
}

func (u *MemoryStorage) GetCity(id int64) (string, error) {
	return u.data[id].City, nil
}

func (u *MemoryStorage) GetLocation(id int64) (string, string, error) {
	return u.data[id].Lat, u.data[id].Lon, nil
}

func (u *MemoryStorage) GetLast(id int64) (string, error) {
	return u.data[id].Last, nil
}

func (u *MemoryStorage) Exists(id int64) (bool, error) {
	_, e := u.data[id]
	return e, nil
}
