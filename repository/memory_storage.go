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

func (u *MemoryStorage) SetSystem(id int64, system bool) {
	currentData := u.data[id]
	currentData.Metric = system
	u.data[id] = currentData
}

func (u *MemoryStorage) SetCity(id int64, city string) {
	currentData := u.data[id]
	currentData.City = city
	u.data[id] = currentData
}

func (u *MemoryStorage) SetLocation(id int64, lat, lon string) {
	currentData := u.data[id]
	currentData.Lat = lat
	currentData.Lon = lon
	u.data[id] = currentData
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

func (u *MemoryStorage) GetCity(id int64) string {
	return u.data[id].City
}

func (u *MemoryStorage) GetLat(id int64) string {
	return u.data[id].Lat
}

func (u *MemoryStorage) GetLon(id int64) string {
	return u.data[id].Lon
}

func (u *MemoryStorage) GetLast(id int64) (weatherCommand string, err error) {
	return u.data[id].Last, nil
}

func (u *MemoryStorage) Exists(id int64) bool {
	_, e := u.data[id]
	return e
}
