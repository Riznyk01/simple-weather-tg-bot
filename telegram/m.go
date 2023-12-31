package telegram

import (
	"SimpleWeatherTgBot/lib/e"
	"SimpleWeatherTgBot/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// Sends a message with text and inline keyboard for weather type selection
func SendMessageWithInlineKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string, buttons ...string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	var inlineButtons []tgbotapi.InlineKeyboardButton
	for _, button := range buttons {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(inlineButtons...))
	msg.ReplyMarkup = inlineKeyboard
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	return err
}

// Sends inline keyboard when requests for weather type is reported after the location is sent.
func SendLocationOptions(bot *tgbotapi.BotAPI, chatID int64, latStr, lonStr string) error {
	chooseWeatherType := fmt.Sprintf("Your location: %s, %v\n%s", latStr, lonStr, types.ChooseOptionMessage)
	err := SendMessageWithInlineKeyboard(bot, chatID, chooseWeatherType, types.CommandCurrentLocation, types.CommandForecastLocation)
	if err != nil {
		return err
	}
	return nil
}

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(e.Wrap("", err))
	}
}
