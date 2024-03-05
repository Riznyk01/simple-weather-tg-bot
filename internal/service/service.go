package service

import (
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/weather_client"
	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	HandleCommand(message *tgbotapi.Message, fname string) (model.UserMessage, error)
	HandleStartCommand(message *tgbotapi.Message, fname string) (model.UserMessage, error)
	HandleHelpCommand() (model.UserMessage, error)
	HandleUnitsCommand(message *tgbotapi.Message) (model.UserMessage, error)
	HandleText(message *tgbotapi.Message) (model.UserMessage, error)
	HandleLocation(message *tgbotapi.Message) (model.UserMessage, error)
	HandleCallbackQuery(callback *tgbotapi.CallbackQuery) (model.UserMessage, error)
	HandleCallbackLast(callback *tgbotapi.CallbackQuery, fname string) (model.UserMessage, error)
}

type Service struct {
	Handler
}

func NewService(log *logr.Logger, repo *repository.Repository, client weather_client.Client) *Service {
	return &Service{
		Handler: NewCommandsHandlerService(log, *repo, client),
	}
}
