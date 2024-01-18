package repository

type WeatherControl interface {
	SetSystem(id int64, system bool) error
	SetCity(id int64, city string) error
	SetLocation(id int64, lat, lon string) error
	SetLast(id int64, last string) error
	GetSystem(id int64) (bool, error)
	GetCity(id int64) (string, error)
	GetLocation(id int64) (string, string, error)
	GetLast(id int64) (string, error)
	AddRequestsCount(id int64) (int, error)
	SetRepliedUserId(chatId, replId int64) error
	GetRepliedUserId(userId int64) (int64, error)
}

type UserControl interface {
	BanUser(id int64) error
	UnbanUser(id int64) error
}

type Repository struct {
	WeatherControl
	UserControl
}

func NewRepository(memoryStor *MemoryStorage) *Repository {
	return &Repository{
		WeatherControl: NewWeatherControlMemoryStorage(memoryStor),
		UserControl:    NewUserControlMemoryStorage(memoryStor),
	}
}
