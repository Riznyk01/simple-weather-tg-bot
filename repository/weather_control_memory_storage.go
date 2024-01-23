package repository

type WeatherControlMemoryStorage struct {
	memoryStor *MemoryStorage
}

func NewWeatherControlMemoryStorage(memoryStor *MemoryStorage) *WeatherControlMemoryStorage {
	return &WeatherControlMemoryStorage{
		memoryStor: memoryStor,
	}
}

// SetSystem sets user's system of measurement.
func (u *WeatherControlMemoryStorage) SetSystem(id int64, system bool) error {
	currentData := u.memoryStor.data[id]
	currentData.Metric = system
	u.memoryStor.data[id] = currentData
	return nil
}

// SetCity gets user's last city.
func (u *WeatherControlMemoryStorage) SetCity(id int64, city string) error {
	currentData := u.memoryStor.data[id]
	currentData.City = city
	u.memoryStor.data[id] = currentData
	return nil
}

// SetLocation gets user's last location.
func (u *WeatherControlMemoryStorage) SetLocation(id int64, lat, lon string) error {
	currentData := u.memoryStor.data[id]
	currentData.Lat = lat
	currentData.Lon = lon
	u.memoryStor.data[id] = currentData
	return nil
}

// SetLast sets last users forecast type.
func (u *WeatherControlMemoryStorage) SetLast(id int64, last string) error {
	currentData := u.memoryStor.data[id]
	currentData.Last = last
	u.memoryStor.data[id] = currentData
	return nil
}

// GetSystem gets user's system of measurement.
func (u *WeatherControlMemoryStorage) GetSystem(id int64) (bool, error) {
	data, _ := u.memoryStor.data[id]
	return data.Metric, nil
}

// GetCity gets user's last city.
func (u *WeatherControlMemoryStorage) GetCity(id int64) (string, error) {
	data, _ := u.memoryStor.data[id]
	if data.City == "" {
		return "", ErrItemIsEmpty
	}
	return data.City, nil
}

// GetLocation gets user's last location.
func (u *WeatherControlMemoryStorage) GetLocation(id int64) (string, string, error) {
	data, _ := u.memoryStor.data[id]
	if data.Lat == "" && data.Lon == "" {
		return "", "", ErrItemIsEmpty
	}
	return data.Lat, data.Lon, nil
}

// GetLast gets user's last weather forecast.
func (u *WeatherControlMemoryStorage) GetLast(id int64) (string, error) {
	data, _ := u.memoryStor.data[id]
	if data.Last == "" {
		return "", ErrItemIsEmpty
	}
	return data.Last, nil
}
