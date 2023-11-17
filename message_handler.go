package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// sendMessageWithKeyboard sends a message with the specified text and keyboard buttons.
func sendMessageWithKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string, buttons ...string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	keyboardButtons := make([]tgbotapi.KeyboardButton, len(buttons))
	for i, button := range buttons {
		keyboardButtons[i] = tgbotapi.NewKeyboardButton(button)
	}
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(keyboardButtons...))
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	return err
}

// sendLocationOptions sends a message with location-related options.
func sendLocationOptions(bot *tgbotapi.BotAPI, chatID int64, latStr, lonStr string) error {
	chooseWeatherType := fmt.Sprintf("Your location: %s, %v\n%s", latStr, lonStr, ChooseOptionMessage)
	return sendMessageWithKeyboard(bot, chatID, chooseWeatherType, CommandForecastLocation, CommandCurrentLocation)
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err := bot.Send(msg)
	if err != nil {
		HandleError("", err)
	}
}
