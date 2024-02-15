package config

import (
	"io"
	"os"
)

type Config struct {
	BotToken      string
	WToken        string
	LogLevel      string
	LogType       string
	BotDebug      bool
	WeatherApiUrl string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig() (botCfg *Config, err error) {

	var weatherToken, botToken string
	if os.Getenv("WEATHER_TOKEN_FILE") != "" {
		weatherToken, err = readSecret(os.Getenv("WEATHER_TOKEN_FILE"))
		if err != nil {
			return nil, err
		}
	} else {
		weatherToken = os.Getenv("WEATHER_TOKEN")
	}

	if os.Getenv("BOT_TOKEN_FILE") != "" {
		botToken, err = readSecret(os.Getenv("BOT_TOKEN_FILE"))
		if err != nil {
			return nil, err
		}
	} else {
		botToken = os.Getenv("BOT_TOKEN")
	}

	return &Config{
		BotToken:      botToken,
		WToken:        weatherToken,
		LogLevel:      os.Getenv("LOG_LEVEL"),
		LogType:       os.Getenv("TYPE_OF_LOG"),
		BotDebug:      DebugStrToBool(os.Getenv("BOT_DEBUG")),
		WeatherApiUrl: os.Getenv("WEATHER_API_URL"),
	}, nil
}

func DebugStrToBool(envDebugVar string) bool {
	return envDebugVar == "true"
}

func NewConfigPostgres() (postgresCfg PostgresConfig, err error) {
	var dbPass string

	if os.Getenv("DB_PASSWORD_FILE") != "" {
		dbPass, err = readSecret(os.Getenv("DB_PASSWORD_FILE"))
		if err != nil {
			return PostgresConfig{}, err
		}
	} else {
		dbPass = os.Getenv("DB_PASSWORD")
	}
	postgresCfg = PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: dbPass,
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	return postgresCfg, nil
}

func readSecret(path string) (secret string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		secret = string(data[:n])
	}
	return secret, nil
}
