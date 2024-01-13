package storage

type Storage interface {
	SetSystem(id int64, system bool)
	SetCity(id int64, city string)
	SetLocation(id int64, lat, lon string)
	SetLast(id int64, last string)
	GetSystem(id int64) (system bool)
	GetCity(id int64) string
	GetLat(id int64) string
	GetLon(id int64) string
	GetLast(id int64) string
	Exists(id int64) bool
}
