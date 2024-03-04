package telegram

import (
	"SimpleWeatherTgBot/internal/model"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Sends a message with text and inline keyboard for service type selection
func (b *Bot) SendMessageWithInlineKeyboard(chatID int64, text string, buttons ...string) {
	//fc := "SendMessageWithInlineKeyboard"
	msg := tgbotapi.NewMessage(chatID, text)

	var inlineButtons []tgbotapi.InlineKeyboardButton
	for _, button := range buttons {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(inlineButtons...))
	msg.ReplyMarkup = inlineKeyboard
	msg.ParseMode = "HTML"
	_, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error(err, "Error occurred while sending message with inline keyb. to the user.")
	}
}

// Sends inline keyboard when requests for service type is reported after the location is sent.
func (b *Bot) SendLocationOptions(chatID int64, latStr, lonStr string) {
	chooseWeatherType := fmt.Sprintf("Your location: %s, %v\n%s", latStr, lonStr, model.MessageChooseOption)
	b.SendMessageWithInlineKeyboard(chatID, chooseWeatherType, model.CallbackCurrentLocation, model.CallbackForecastLocation)
}

func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error(err, "Error occurred while sending message to the user.")
	}
}
