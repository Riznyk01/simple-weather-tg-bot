package main

import (
	"SimpleWeatherTgBot/lib/e"
	"SimpleWeatherTgBot/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// sends a message with text and buttons for weather type selection
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

// sends when requests for weather type is reported after the location is sent
func sendLocationOptions(bot *tgbotapi.BotAPI, chatID int64, latStr, lonStr string) error {
	chooseWeatherType := fmt.Sprintf("Your location: %s, %v\n%s", latStr, lonStr, types.ChooseOptionMessage)
	return sendMessageWithKeyboard(bot, chatID, chooseWeatherType, types.CommandForecastLocation, types.CommandCurrentLocation)
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(e.Wrap("", err))
	}
}
