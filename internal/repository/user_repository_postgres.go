package repository

import (
	"SimpleWeatherTgBot/internal/model"
	"database/sql"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserRepositoryPostgres struct {
	log *logr.Logger
	db  *sqlx.DB
}

type UserDataPostgres struct {
	City            sql.NullString
	Lat             sql.NullString
	Lon             sql.NullString
	Metric          bool
	LastWeatherType sql.NullString
}

func NewUserRepository(log *logr.Logger, db *sqlx.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		log: log,
		db:  db,
	}
}

// SetUserMeasurementSystem sets user's system of measurement.
func (r *UserRepositoryPostgres) SetUserMeasurementSystem(id int64, system bool) error {
	//fc := "SetUserMeasurementSystem"

	q := fmt.Sprintf("UPDATE %s SET metric = $1 WHERE id = $2", usersTable)
	_, err := r.db.Exec(q, system, id)
	if err != nil {
		r.log.Error(err, "Error updating user's preferred system of measurement")
		return err
	}

	return nil
}

// AddUsersSchedule ...
func (r *UserRepositoryPostgres) AddUsersSchedule(id int64, scheduleCity string, scheduleTime time.Time, weatherType string, timezoneOffset float64) error {
	q := fmt.Sprintf("INSERT INTO %s (id, city, schedule_time, weather_type, timezone_offset) VALUES ($1, $2, $3, $4, $5)", schedulesTable)
	_, err := r.db.Exec(q, id, scheduleCity, scheduleTime, weatherType, timezoneOffset)
	if err != nil {
		r.log.Error(err, "Error adding user's schedule")
		return err
	}

	return nil
}

// DeleteUsersSchedule ...
func (r *UserRepositoryPostgres) DeleteUsersSchedule(id int64, scheduleCity string, scheduleTime time.Time, weatherType string, timezoneOffset float64) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND city = $2 AND schedule_time = $3 AND weather_type = $4 AND timezone_offset = $5", schedulesTable)
	_, err := r.db.Exec(q, id, scheduleCity, scheduleTime, weatherType, timezoneOffset)
	if err != nil {
		r.log.Error(err, "Error deleting user's schedule")
		return err
	}

	return nil
}

// SetUserLastInputCity sets the user's last input city for weather forecast.
func (r *UserRepositoryPostgres) SetUserLastInputCity(id int64, city string) error {
	//fc := "SetUserLastInputCity"

	q := fmt.Sprintf("UPDATE %s SET city = $1 WHERE id = $2", usersTable)
	_, err := r.db.Exec(q, city, id)
	if err != nil {
		r.log.Error(err, "Error updating user's preferred city")
		return err
	}

	return nil
}

// SetUserLastInputLocation sets the user's last input location for weather forecast.
func (r *UserRepositoryPostgres) SetUserLastInputLocation(id int64, lat, lon string) error {
	//fc := "SetUserLastInputLocation"

	q := fmt.Sprintf("UPDATE %s SET lat = $1, lon = $2 WHERE id = $3", usersTable)
	_, err := r.db.Exec(q, lat, lon, id)
	if err != nil {
		r.log.Error(err, "Error updating user's preferred location")
		return err
	}
	return nil
}

// SetUserLastWeatherCommand sets the user's last input forecast type.
func (r *UserRepositoryPostgres) SetUserLastWeatherCommand(userId int64, last string) error {
	//fc := "SetUserLastWeatherCommand"

	q := fmt.Sprintf("UPDATE %s SET last = $1 WHERE id = $2", usersTable)
	_, err := r.db.Exec(q, last, userId)
	if err != nil {
		r.log.Error(err, "Error updating last weather command")
		return err
	}
	return nil
}

// GetUserById gets the user's data from the database.
func (r *UserRepositoryPostgres) GetUserById(userId int64) (model.UserData, error) {
	//fc := "GetUserById"

	q := fmt.Sprintf("SELECT city, lat, lon, metric, last FROM %s WHERE id = $1", usersTable)
	row := r.db.QueryRow(q, userId)

	var user UserDataPostgres
	err := row.Scan(&user.City, &user.Lat, &user.Lon, &user.Metric, &user.LastWeatherType)
	if err != nil {
		r.log.Error(err, "Error fetching user from the database")
		return model.UserData{}, err
	}

	userData := model.UserData{
		City:            handleNullString(user.City),
		Lat:             handleNullString(user.Lat),
		Lon:             handleNullString(user.Lon),
		Metric:          user.Metric,
		LastWeatherType: handleNullString(user.LastWeatherType),
	}

	return userData, nil
}

func handleNullString(nullStr sql.NullString) string {
	if nullStr.Valid {
		return nullStr.String
	}
	return ""
}

// CreateUserById creates a user in the database.
func (r *UserRepositoryPostgres) CreateUserById(userId int64) error {
	//fc := "CreateUserById"

	q := fmt.Sprintf("INSERT INTO %s (id, metric) VALUES ($1, true)", usersTable)
	_, err := r.db.Exec(q, userId)
	if err != nil {
		r.log.Error(err, "Error inserting user into the database")
		return err
	}
	return nil
}
