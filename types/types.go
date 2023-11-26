package types

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type WeatherResponse5d3h struct {
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`
		Visibility int     `json:"visibility"`
		Pop        float64 `json:"pop"`
		Sys        struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
		Rain  struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain,omitempty"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

type UserData struct {
	City   string
	Lat    string
	Lon    string
	Metric bool
	Last   string
}

type MemoryStorage struct {
	Data map[int64]UserData
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Data: make(map[int64]UserData),
	}
}

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

func (u *MemoryStorage) SetSystem(id int64, system bool) {
	currentData := u.Data[id]
	currentData.Metric = system
	u.Data[id] = currentData
}

func (u *MemoryStorage) SetCity(id int64, city string) {
	currentData := u.Data[id]
	currentData.City = city
	u.Data[id] = currentData
}

func (u *MemoryStorage) SetLocation(id int64, lat, lon string) {
	currentData := u.Data[id]
	currentData.Lat = lat
	currentData.Lon = lon
	u.Data[id] = currentData
}

func (u *MemoryStorage) SetLast(id int64, last string) {
	currentData := u.Data[id]
	currentData.Last = last
	u.Data[id] = currentData
}

func (u *MemoryStorage) GetSystem(id int64) bool {
	return u.Data[id].Metric
}

func (u *MemoryStorage) GetCity(id int64) string {
	return u.Data[id].City
}

func (u *MemoryStorage) GetLat(id int64) string {
	return u.Data[id].Lat
}

func (u *MemoryStorage) GetLon(id int64) string {
	return u.Data[id].Lon
}

func (u *MemoryStorage) GetLast(id int64) string {
	return u.Data[id].Last
}

func (u *MemoryStorage) Exists(id int64) bool {
	_, e := u.Data[id]
	return e
}
