package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func SetupLogger() *logrus.Logger {
	log := logrus.New()
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)
	setLogType(log)
	return log
}
func setLogType(log *logrus.Logger) {
	switch os.Getenv("TYPE_OF_LOG") {
	case "TEXTLOG":
		log.SetFormatter(&logrus.TextFormatter{})
	case "JSONLOG":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.JSONFormatter{})
	}
}
