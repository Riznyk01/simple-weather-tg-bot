package repository

import (
	"SimpleWeatherTgBot/internal/model"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepositoryPostgres struct {
	log *logrus.Logger
	db  *sqlx.DB
}

type UserDataPostgres struct {
	City   sql.NullString
	Lat    sql.NullString
	Lon    sql.NullString
	Metric bool
	Last   sql.NullString
}

func NewUserRepository(log *logrus.Logger, db *sqlx.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		log: log,
		db:  db,
	}
}
func (r *UserRepositoryPostgres) SetSystem(id int64, system bool) error {
	fc := "SetSystem"

	q := fmt.Sprintf("UPDATE user_data SET metric = $1 WHERE id = $2")
	_, err := r.db.Exec(q, system, id)
	if err != nil {
		r.log.Errorf("%s: Error updating system: %v", fc, err)
		return err
	}

	return nil
}

func (r *UserRepositoryPostgres) SetCity(id int64, city string) error {
	fc := "SetCity"

	q := fmt.Sprintf("UPDATE user_data SET city = $1 WHERE id = $2")
	_, err := r.db.Exec(q, city, id)
	if err != nil {
		r.log.Errorf("%s: Error updating city: %v", fc, err)
		return err
	}

	return nil
}

func (r *UserRepositoryPostgres) SetLocation(id int64, lat, lon string) error {
	fc := "SetLocation"

	q := fmt.Sprintf("UPDATE user_data SET lat = $1, lon = $2 WHERE id = $3")
	_, err := r.db.Exec(q, lat, lon, id)
	if err != nil {
		r.log.Errorf("%s: Error updating location: %v", fc, err)
		return err
	}
	return nil
}

func (r *UserRepositoryPostgres) SetLastWeatherCommand(userId int64, last string) error {
	fc := "SetLastWeatherCommand"

	q := fmt.Sprintf("UPDATE %s SET last = $1 WHERE id = $2", usersTable)
	_, err := r.db.Exec(q, last, userId)
	if err != nil {
		r.log.Errorf("%s: Error updating last weather command: %v", fc, err)
		return err
	}
	return nil
}

func (r *UserRepositoryPostgres) GetUserById(userId int64) (model.UserData, error) {
	fc := "GetUserById"

	q := fmt.Sprintf("SELECT city, lat, lon, metric, last FROM %s WHERE id = $1", usersTable)
	row := r.db.QueryRow(q, userId)

	var user UserDataPostgres
	err := row.Scan(&user.City, &user.Lat, &user.Lon, &user.Metric, &user.Last)
	if err != nil {
		r.log.Errorf("%s: Error fetching user: %v", fc, err)
		return model.UserData{}, err
	}

	userData := model.UserData{
		City:   handleNullString(user.City),
		Lat:    handleNullString(user.Lat),
		Lon:    handleNullString(user.Lon),
		Metric: user.Metric,
		Last:   handleNullString(user.Last),
	}

	return userData, nil
}

func handleNullString(nullStr sql.NullString) string {
	if nullStr.Valid {
		return nullStr.String
	}
	return ""
}

func (r *UserRepositoryPostgres) CreateUser(userId int64) error {
	fc := "CreateUser"

	q := fmt.Sprintf("INSERT INTO %s (id, metric) VALUES ($1, true)", usersTable)
	_, err := r.db.Exec(q, userId)
	if err != nil {
		r.log.Errorf("%s: Error inserting user: %v", fc, err)
		return err
	}
	return nil
}
