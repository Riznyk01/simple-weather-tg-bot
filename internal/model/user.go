package model

type UserData struct {
	City            string
	Lat             string
	Lon             string
	Metric          bool
	LastWeatherType string
}

type UserMessage struct {
	Text    string
	Buttons []string
}
