package model

import "time"

type ScheduleData struct {
	ID             int64     `db:"id"`
	City           string    `db:"city"`
	ScheduleTime   time.Time `db:"schedule_time"`
	WeatherType    string    `db:"weather_type"`
	TimezoneOffset float64   `db:"timezone_offset"`
	Units          bool      `db:"units"`
}
