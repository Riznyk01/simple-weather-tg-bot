package repository

type WeatherUserControl interface {
	SetSystem(id int64, system bool)
	SetCity(id int64, city string)
	SetLocation(id int64, lat, lon string)
	SetLast(id int64, last string) error
	GetSystem(id int64) (bool, error)
	GetCity(id int64) string
	GetLat(id int64) string
	GetLon(id int64) string
	GetLast(id int64) (weatherCommand string, err error)
	Exists(id int64) bool
}

type Repository struct {
	WeatherUserControl
}

func NewRepository() *Repository {
	return &Repository{
		WeatherUserControl: NewMemoryStorage(),
	}
}
