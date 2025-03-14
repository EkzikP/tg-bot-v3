package menus

import (
	"fmt"
	"github.com/EkzikP/tg-bot-v3/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MenuBuilder struct {
	MainMenu     []models.Menu
	MyAlarmMenu  []models.Menu
	RequestPhone models.RequestPhone
}

func New() *MenuBuilder {
	mainMenu := []models.Menu{
		{"Получить информацию по объекту", "GetInfoObject"},
		{"Получить список ответственных лиц объекта", "GetCustomers"},
		{"Проверка КТС", "ChecksKTS"},
		{"Управление доступом в MyAlarm", "MyAlarm"},
		{"Получить список разделов", "GetParts"},
		{"Получить список шлейфов", "GetZones"},
		{"Завершить работу с объектом", "Finish"},
	}
	myAlarmMenu := []models.Menu{
		{"Список пользователей MyAlarm объекта", "GetUsersMyAlarm"},
		{"Список объектов пользователя MyAlarm", "GetUserObjectMyAlarm"},
		{"Забрать доступ к MyAlarm", "PutDelUserMyAlarm"},
		{"Предоставить доступ к MyAlarm", "PutAddUserMyAlarm"},
		{"Модифицировать виртуальную КТС", "PutChangeVirtualKTS"},
		{"Назад", "Back"},
		{"Завершить работу с объектом", "Finish"},
	}
	requestPhoneMenu := models.RequestPhone{
		Text:           "Отправить номер телефона",
		RequestContact: true,
	}
	return &MenuBuilder{
		MainMenu:     mainMenu,
		MyAlarmMenu:  myAlarmMenu,
		RequestPhone: requestPhoneMenu,
	}
}

func (b *MenuBuilder) BuildMainMenu(chatID int64, numberObject string) tgbotapi.MessageConfig {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, button := range b.MainMenu {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(button.Text, button.CallbackData)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	text := fmt.Sprintf("Работа с объектом %s!\nВыберите пункт меню:", numberObject)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = &keyboard
	return msg
}

func (b *MenuBuilder) BuildMyAlarmMenu(chatID int64, numberObject string) tgbotapi.MessageConfig {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, button := range b.MyAlarmMenu {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(button.Text, button.CallbackData)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	text := fmt.Sprintf("Работа с объектом %s!\nВыберите пункт меню:", numberObject)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = &keyboard
	return msg
}

func (b *MenuBuilder) RequestContact(chatID int64) tgbotapi.MessageConfig {

	button := tgbotapi.KeyboardButton{Text: b.RequestPhone.Text, RequestContact: b.RequestPhone.RequestContact}
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        [][]tgbotapi.KeyboardButton{{button}},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		Selective:       true,
	}

	msg := tgbotapi.NewMessage(chatID, "Отправьте ваш номер телефона, нажав на кнопку ниже.")
	msg.ReplyMarkup = &keyboard

	return msg
}

func BackAndFinish() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	btn := tgbotapi.NewInlineKeyboardButtonData("Назад", "Back")
	var row []tgbotapi.InlineKeyboardButton
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	btnFinish := tgbotapi.NewInlineKeyboardButtonData("Завершить работу с объектом", "Finish")
	var rowFinish []tgbotapi.InlineKeyboardButton
	rowFinish = append(rowFinish, btnFinish)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowFinish)

	return keyboard
}

func CheckKTS(CheckPanicId string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	btn := tgbotapi.NewInlineKeyboardButtonData("Получить результат проверки КТС", "ResultCheckKTS")
	var row []tgbotapi.InlineKeyboardButton
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, BackAndFinish().InlineKeyboard...)
	return keyboard
}
