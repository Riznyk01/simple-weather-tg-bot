package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type RespMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type Geocoding struct {
	//Name string `json:"name"`
	//LocalNames struct {
	//	Iu          string `json:"iu,omitempty"`
	//} `json:"local_names,omitempty"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	//Country string  `json:"country"`
	//State   string  `json:"state"`
}
