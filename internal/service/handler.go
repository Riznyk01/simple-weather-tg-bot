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
	"strconv"
	"strings"
	"time"
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
	} else if strings.HasPrefix(message.Text, "/add") {
		return h.HandleAddSchedule(message)
	} else if strings.HasPrefix(message.Text, text.CommandDeleteSchedule) {
		return h.HandleDeleteSchedule(message)
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

// HandleAddSchedule ...
func (h *CommandsHandlerService) HandleAddSchedule(message *tgbotapi.Message) (model.UserMessage, error) {
	parts := strings.Split(message.Text, "_")
	if len(parts) != 6 {
		return model.UserMessage{Text: "please check if you typed the correct command, like /add_18:00_2_cityname_weathertype_metricunits", Buttons: nil}, nil
	}

	timeStr := parts[1]
	offsetStr := parts[2]
	scheduleCity := parts[3]
	weatherType := parts[4]

	var metricUnits bool
	if parts[5] == "true" {
		metricUnits = true
	}

	t, _ := time.Parse("15:04", timeStr)

	offset, err := strconv.ParseFloat(offsetStr, 64)
	if err != nil {
		return model.UserMessage{Text: "invalid timezone offset", Buttons: nil}, nil
	}

	// Convert local time to UTC
	offsetHours := int(offset)
	offsetMinutes := int((offset - float64(offsetHours)) * 60)
	loc := time.FixedZone("UserTZ", offsetHours*3600+offsetMinutes*60)

	// Create local time considering the offset
	localTime := time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, loc)
	utcTime := localTime.UTC()

	err = h.repo.AddUsersSchedule(message.Chat.ID, scheduleCity, utcTime, weatherType, offset, metricUnits)
	if err != nil {
		return model.UserMessage{Text: err.Error(), Buttons: nil}, err
	}
	return model.UserMessage{Text: "the schedule was added to the list", Buttons: nil}, nil
}

// HandleDeleteSchedule ...
func (h *CommandsHandlerService) HandleDeleteSchedule(message *tgbotapi.Message) (model.UserMessage, error) {
	parts := strings.Split(message.Text, "_")
	if len(parts) != 2 {
		return model.UserMessage{Text: fmt.Sprintf("please check if you typed the correct command, like %s", text.CommandDeleteSchedule),
			Buttons: nil}, nil
	}

	scheduleCity := parts[1]

	err := h.repo.DeleteUsersSchedule(message.Chat.ID, scheduleCity)
	if err != nil {
		return model.UserMessage{Text: err.Error(), Buttons: nil}, err
	}
	return model.UserMessage{Text: "the schedule was deleted from the list", Buttons: nil}, nil
}

// HandleRemoveSchedule ...
func (h *CommandsHandlerService) HandleRemoveSchedule() (model.UserMessage, error) {
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
				h.repo.IncrementUserUsageCount(callback.Message.Chat.ID)
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
		if user.LastWeatherType == "" {
			return model.UserMessage{Text: fmt.Sprintf(text.MsgLastDataUnavailable, fname), Buttons: nil}, err
		} else {
			userMessage, err := h.client.GetWeatherForecast(user)
			if err != nil {
				return model.UserMessage{Text: userMessage, Buttons: nil}, err
			} else {
				h.repo.IncrementUserUsageCount(callback.Message.Chat.ID)
				return model.UserMessage{Text: userMessage, Buttons: []string{text.CallbackLast}}, err
			}
		}
	}
}

// HandleSchedule ...
func (h *CommandsHandlerService) HandleSchedule() (int64, model.UserMessage, error) {

	schedules, err := h.repo.GetSchedulesByCurrentTime()
	if err != nil {
		h.log.Error(err, "Error fetching schedules")
		return 0, model.UserMessage{}, err
	}
	//TODO: add metric fetching
	if len(schedules) > 0 {
		for _, schedule := range schedules {
			time.Sleep(1 * time.Second)
			userMessage, err := h.client.GetWeatherForecast(model.UserData{
				schedule.City,
				"",
				"",
				true,
				schedule.WeatherType})
			if err != nil {
				return 0, model.UserMessage{Text: userMessage, Buttons: nil}, err
			} else {
				return schedule.ID, model.UserMessage{Text: userMessage, Buttons: []string{text.CallbackLast}}, err
			}
		}
	}
	return 0, model.UserMessage{Text: "", Buttons: nil}, nil
}
