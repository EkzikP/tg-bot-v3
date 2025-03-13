package menus

import (
	"fmt"
	"github.com/EkzikP/tg-bot-v3/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MenuBuilder struct {
	MainMenu     []models.Menu
	MyAlarmMenu  []models.Menu
	BackMenu     []models.Menu
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
	backMenu := []models.Menu{
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
		BackMenu:     backMenu,
		RequestPhone: requestPhoneMenu,
	}
}

func (b *MenuBuilder) BuildMainMenu(chatID int64, numberObject string) tgbotapi.MessageConfig {
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
