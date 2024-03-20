package service

import (
	"SimpleWeatherTgBot/internal/model"
	"SimpleWeatherTgBot/internal/repository"
	"SimpleWeatherTgBot/internal/text"
	"SimpleWeatherTgBot/internal/weather_client"
	"fmt"
	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lib/pq"
)

type CommandsHandlerService struct {
	log    *logr.Logger
	repo   repository.Repository
	client weather_client.Client
}

func NewCommandsHandlerService(log *logr.Logger, repo repository.Repository, client weather_client.Client) *CommandsHandlerService {
	return &CommandsHandlerService{
		log:    log,
		repo:   repo,
		client: client,
	}
}

// HandleCommand processes commands from the user.
func (h *CommandsHandlerService) HandleCommand(message *tgbotapi.Message, fname string) (model.UserMessage, error) {
	if message.Text == text.CommandMetricUnits || message.Text == text.CommandNonMetricUnits {
		return h.HandleUnitsCommand(message)
	} else if message.Text == text.CommandStart {
		return h.HandleStartCommand(message, fname)
	} else if message.Text == text.CommandHelp {
		return h.HandleHelpCommand()
	}
	//change
	return model.UserMessage{}, nil
}

// HandleStartCommand handles the /start command.
func (h *CommandsHandlerService) HandleStartCommand(message *tgbotapi.Message, fname string) (model.UserMessage, error) {
	err := h.repo.CreateUserById(message.Chat.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			// Checking if the error is a unique constraint violation
			if pgErr.Code == "23505" {
				return model.UserMessage{Text: fmt.Sprintf("%s\n%s", text.MsgAlreadyStarted, text.MsgHelp), Buttons: nil}, nil
			}
		}
		return model.UserMessage{Text: text.ErrWhileExecuting, Buttons: nil}, err
	} else {
		return model.UserMessage{Text: fmt.Sprintf(text.MsgWelcome, fname) + text.MsgHelp, Buttons: nil}, nil
	}
}

// HandleHelpCommand handles the /help command.
func (h *CommandsHandlerService) HandleHelpCommand() (model.UserMessage, error) {
	return model.UserMessage{Text: text.MsgHelp, Buttons: nil}, nil
}

// HandleUnitsCommand handles the /metric and /non-metric commands.
func (h *CommandsHandlerService) HandleUnitsCommand(message *tgbotapi.Message) (model.UserMessage, error) {
	m := false

	if message.Text == text.CommandMetricUnits {
		m = true
	}
	err := h.repo.SetUserMeasurementSystem(message.Chat.ID, m)
	if err != nil {
		return model.UserMessage{Text: text.MsgSetUsersSystemError, Buttons: nil}, err
	} else {
		return model.UserMessage{Text: text.MsgMetricUnitChanged, Buttons: nil}, nil
	}
}

// HandleText processes text from the user.
func (h *CommandsHandlerService) HandleText(message *tgbotapi.Message) (model.UserMessage, error) {
	if !containsEmoji(message.Text) {
		err := h.repo.SetUserLastInputCity(message.Chat.ID, message.Text)
		if err != nil {
			return model.UserMessage{Text: text.MsgSetUsersCityError, Buttons: nil}, nil
		} else {
			return model.UserMessage{Text: text.MsgChooseOption, Buttons: []string{text.CallbackCurrent, text.CallbackForecast, text.CallbackToday}}, nil
		}
	} else {
		return model.UserMessage{Text: text.MsgUnsupportedMessageType, Buttons: nil}, nil
	}
}

// containsEmoji returns true if the text contains emojis.
func containsEmoji(text string) bool {
	for _, char := range text {
		if char >= '\U0001F600' && char <= '\U0001F64F' {
			return true
		}
	}
	return false
}

// HandleLocation processes location messages from the user.
func (h *CommandsHandlerService) HandleLocation(message *tgbotapi.Message) (model.UserMessage, error) {
	uLat, uLon := fmt.Sprintf("%f", message.Location.Latitude), fmt.Sprintf("%f", message.Location.Longitude)
	err := h.repo.SetUserLastInputLocation(message.Chat.ID, uLat, uLon)
	if err != nil {
		return model.UserMessage{Text: text.MsgSetUsersLocationError, Buttons: nil}, err
	} else {
		return model.UserMessage{Text: fmt.Sprintf("Your location: %s, %v\n%s", uLat, uLon, text.MsgChooseOption), Buttons: []string{text.CallbackCurrentLocation, text.CallbackForecastLocation, text.CallbackTodayLocation}}, err
	}
}

// HandleCallbackQuery handles callback queries from the user.
func (h *CommandsHandlerService) HandleCallbackQuery(callback *tgbotapi.CallbackQuery) (model.UserMessage, error) {

	_ = h.repo.SetUserLastWeatherCommand(callback.Message.Chat.ID, callback.Data)
	user, err := h.repo.GetUserById(callback.Message.Chat.ID)
	if err != nil {
		return model.UserMessage{Text: text.ErrWhileExecuting, Buttons: nil}, err
	} else {
		userMessage, err := h.client.GetWeatherForecast(user)
		if err != nil {
			return model.UserMessage{Text: userMessage, Buttons: nil}, err
		} else {
			if userMessage == text.TryAnother {
				return model.UserMessage{Text: text.TryAnother, Buttons: nil}, err
			} else {
				return model.UserMessage{Text: userMessage, Buttons: []string{text.CallbackLast}}, err
			}
		}
	}
}

// HandleCallbackLast handles callback queries from the user with the "repeat last" command.
func (h *CommandsHandlerService) HandleCallbackLast(callback *tgbotapi.CallbackQuery, fname string) (model.UserMessage, error) {

	user, err := h.repo.GetUserById(callback.Message.Chat.ID)
	if err != nil {
		return model.UserMessage{Text: text.ErrWhileExecuting, Buttons: nil}, err
	} else {
		if user.Last == "" {
			return model.UserMessage{Text: fmt.Sprintf(text.MsgLastDataUnavailable, fname), Buttons: nil}, err
		} else {
			userMessage, err := h.client.GetWeatherForecast(user)
			if err != nil {
				return model.UserMessage{Text: userMessage, Buttons: nil}, err
			} else {
				return model.UserMessage{Text: userMessage, Buttons: []string{text.CallbackLast}}, err
			}
		}
	}
}
