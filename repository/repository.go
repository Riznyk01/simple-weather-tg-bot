package repository

type WeatherUserControl interface {
	SetSystem(id int64, system bool) error
	SetCity(id int64, city string) error
	SetLocation(id int64, lat, lon string) error
	SetLast(id int64, last string) error
	GetSystem(id int64) (bool, error)
	GetCity(id int64) (string, error)
	GetLocation(id int64) (string, string, error)
	GetLast(id int64) (string, error)
	Exists(id int64) (bool, error)
	AddRequestsCount(id int64) int
}

type Repository struct {
	WeatherUserControl
}

func NewRepository() *Repository {
	return &Repository{
		WeatherUserControl: NewMemoryStorage(),
	}
}
