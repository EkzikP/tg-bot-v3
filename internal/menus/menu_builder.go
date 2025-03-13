package menus

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MenuBuilder struct{}

func New() *MenuBuilder {
	return &MenuBuilder{}
}

func (b *MenuBuilder) MainMenu(chatID int64, numberObject string) tgbotapi.MessageConfig {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Информация", "GetInfoObject"),
			tgbotapi.NewInlineKeyboardButtonData("КТС", "ChecksKTS"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Объект %s\nВыберите действие:", numberObject))
	msg.ReplyMarkup = keyboard
	return msg
}

// Другие методы создания меню
