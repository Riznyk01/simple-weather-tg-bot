package config

import (
	"errors"
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
	botCfg = &Config{}

	if os.Getenv("WEATHER_TOKEN_FILE") != "" {
		botCfg.WToken, err = readSecret(os.Getenv("WEATHER_TOKEN_FILE"))
		if err != nil {
			return nil, err
		}
	} else {
		botCfg.WToken = os.Getenv("WEATHER_TOKEN")
	}

	if os.Getenv("BOT_TOKEN_FILE") != "" {
		botCfg.BotToken, err = readSecret(os.Getenv("BOT_TOKEN_FILE"))
		if err != nil {
			return nil, err
		}
	} else {
		botCfg.BotToken = os.Getenv("BOT_TOKEN")
	}

	botCfg.LogLevel = os.Getenv("LOG_LEVEL")
	botCfg.LogType = os.Getenv("TYPE_OF_LOG")
	botCfg.BotDebug = DebugStrToBool(os.Getenv("BOT_DEBUG"))
	botCfg.WeatherApiUrl = os.Getenv("WEATHER_API_URL")

	if botCfg.LogLevel == "" || botCfg.LogType == "" || botCfg.WeatherApiUrl == "" || botCfg.BotToken == "" || botCfg.WToken == "" {
		return nil, errors.New("some fields in the app config are empty")
	}

	return botCfg, nil
}

func DebugStrToBool(envDebugVar string) bool {
	return envDebugVar == "true"
}

func NewConfigPostgres() (postgresCfg *PostgresConfig, err error) {
	postgresCfg = &PostgresConfig{}
	if os.Getenv("DB_PASSWORD_FILE") != "" {
		postgresCfg.Password, err = readSecret(os.Getenv("DB_PASSWORD_FILE"))
		if err != nil {
			return nil, err
		}
	} else {
		postgresCfg.Password = os.Getenv("DB_PASSWORD")
	}

	postgresCfg.Host = os.Getenv("DB_HOST")
	postgresCfg.Port = os.Getenv("DB_PORT")
	postgresCfg.Username = os.Getenv("DB_USERNAME")
	postgresCfg.DBName = os.Getenv("DB_NAME")
	postgresCfg.SSLMode = os.Getenv("DB_SSLMODE")

	if postgresCfg.Host == "" || postgresCfg.Port == "" || postgresCfg.Username == "" || postgresCfg.DBName == "" || postgresCfg.SSLMode == "" {
		return nil, errors.New("some fields in the PostgreSQL config are empty")
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
