package repository

import (
	"SimpleWeatherTgBot/config"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "user_data"
)

func NewPostgresDB(pgCfg config.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		pgCfg.Host, pgCfg.Port, pgCfg.Username, pgCfg.DBName, pgCfg.Password, pgCfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
