package telegram

import (
	"SimpleWeatherTgBot/internal/text"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendMessageWithInlineKeyboard sends a message with text and inline keyboard for service type selection
func (b *Bot) SendMessageWithInlineKeyboard(chatID int64, msgText string, buttons ...string) {

	msg := tgbotapi.NewMessage(chatID, msgText)

	var inlineButtons []tgbotapi.InlineKeyboardButton
	for _, button := range buttons {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(inlineButtons...))
	msg.ReplyMarkup = inlineKeyboard
	msg.ParseMode = "HTML"
	_, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error(err, text.ErrWhileSendingInline)
	}
}

// SendLocationOptions sends inline keyboard when requests for service type is reported after the location is sent.
func (b *Bot) SendLocationOptions(chatID int64, latStr, lonStr string) {
	chooseWeatherType := fmt.Sprintf("Your location: %s, %v\n%s", latStr, lonStr, text.MsgChooseOption)
	b.SendMessageWithInlineKeyboard(chatID, chooseWeatherType, text.CallbackCurrentLocation, text.CallbackForecastLocation)
}

func (b *Bot) SendMessage(chatID int64, msgText string) {
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = "HTML"
	_, err := b.bot.Send(msg)
	if err != nil {
		b.log.Error(err, text.ErrWhileSendingMsg)
	}
}
