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

// SetSystem sets user's system of measurement.
func (u *MemoryStorage) SetSystem(id int64, system bool) error {
	currentData := u.data[id]
	currentData.Metric = system
	u.data[id] = currentData
	return nil
}

// SetCity gets user's last city.
func (u *MemoryStorage) SetCity(id int64, city string) error {
	currentData := u.data[id]
	currentData.City = city
	u.data[id] = currentData
	return nil
}

// SetLocation gets user's last location.
func (u *MemoryStorage) SetLocation(id int64, lat, lon string) error {
	currentData := u.data[id]
	currentData.Lat = lat
	currentData.Lon = lon
	u.data[id] = currentData
	return nil
}

// SetLast sets last users forecast type.
func (u *MemoryStorage) SetLast(id int64, last string) error {
	currentData := u.data[id]
	currentData.Last = last
	u.data[id] = currentData
	return nil
}

// GetSystem gets user's system of measurement.
func (u *MemoryStorage) GetSystem(id int64) (bool, error) {
	return u.data[id].Metric, nil
}

// GetCity gets user's last city.
func (u *MemoryStorage) GetCity(id int64) (string, error) {
	return u.data[id].City, nil
}

// GetLocation gets user's last location.
func (u *MemoryStorage) GetLocation(id int64) (string, string, error) {
	return u.data[id].Lat, u.data[id].Lon, nil
}

// GetLast gets user's last weather forecast.
func (u *MemoryStorage) GetLast(id int64) (string, error) {
	return u.data[id].Last, nil
}

// Exists checks if the user exists.
func (u *MemoryStorage) Exists(id int64) (bool, error) {
	_, e := u.data[id]
	return e, nil
}
