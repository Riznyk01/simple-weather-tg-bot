package model

type UserData struct {
	City   string
	Lat    string
	Lon    string
	Metric bool
	Last   string
}

type UserMessage struct {
	Text    string
	Buttons []string
}
