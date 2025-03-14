package handlers

import (
	"context"
	"fmt"
	"github.com/EkzikP/tg-bot-v3/internal/menus"
	"github.com/EkzikP/tg-bot-v3/internal/models"
	"github.com/EkzikP/tg-bot-v3/internal/services"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler struct {
	Bot           *tgbotapi.BotAPI
	Andromeda     *services.AndromedaService
	Operations    map[int64]*models.Operation
	PhoneEngineer map[string]string
}

func NewCallbackHandler(bot *tgbotapi.BotAPI, phoneEngineer map[string]string, andromeda *services.AndromedaService) *CallbackHandler {
	return &CallbackHandler{
		Bot:           bot,
		Andromeda:     andromeda,
		Operations:    make(map[int64]*models.Operation),
		PhoneEngineer: phoneEngineer,
	}
}

func (h *CallbackHandler) sendMessage(msg tgbotapi.MessageConfig) {
	h.Bot.Send(msg)
}

func (h *CallbackHandler) HandleCallback(ctx context.Context, update tgbotapi.Update, tgUser *sync.Map) {
	callback := update.CallbackQuery
	chatID := callback.Message.Chat.ID
	currentOperations := h.Operations[chatID]

	switch callback.Data {
	case "Finish":
		h.finish(chatID, currentOperations)
	case "Back":
		h.back(chatID, currentOperations, callback)
	case "GetCustomers":
		currentOperations.Update("CurrentRequest", callback.Data)
		h.getCustomers(chatID, currentOperations)
	case "GetInfoObject":
		h.handleGetInfo(chatID)
	case "ChecksKTS", "ResultCheckKTS":
		currentOperations.Update("CurrentRequest", callback.Data)
		h.handleChecksKTS(chatID, currentOperations)

		// Обработка других callback-ов
	}
}

func (h *CallbackHandler) finish(chatID int64, currentOperations *models.Operation) {
	text := fmt.Sprintf("Завершена работа с объектом %s", currentOperations.NumberObject)
	msg := tgbotapi.NewMessage(chatID, text)
	h.sendMessage(msg)
	unpinMessage := tgbotapi.UnpinAllChatMessagesConfig{
		ChatID: chatID,
	}
	h.Bot.Send(unpinMessage)

	currentOperations = models.New()
	msg = tgbotapi.NewMessage(chatID, "Введите пультовый номер объекта!")
	h.sendMessage(msg)
}

func (h *CallbackHandler) back(chatID int64, currentOperations *models.Operation, callback *tgbotapi.CallbackQuery) {
	mb := menus.New()
	text := fmt.Sprintf("Работа с объектом %s", currentOperations.NumberObject)
	msg := tgbotapi.MessageConfig{}
	if currentOperations.CurrentMenu == "MyAlarmMenu" && currentOperations.CurrentRequest == "MyAlarm" {
		text += "\nПодменю MyAlarm"
		msg = tgbotapi.NewMessage(chatID, text)
		currentOperations.Update("CurrentRequest", "")
		currentOperations.Update("CurrentMenu", "MainMenu")
		msg = mb.BuildMainMenu(chatID, currentOperations.NumberObject)
		h.sendMessage(msg)
		return
	} else if currentOperations.CurrentMenu == "MyAlarmMenu" {
		msg = tgbotapi.NewMessage(chatID, text)
		text += "\nПодменю MyAlarm"
		msg = tgbotapi.NewMessage(chatID, text)
		currentOperations.Update("CurrentRequest", "MyAlarm")
		msg = mb.BuildMyAlarmMenu(chatID, currentOperations.NumberObject)
		h.sendMessage(msg)
		return
	}

	msg = tgbotapi.NewMessage(chatID, text)
	currentOperations.Update("CurrentRequest", "")
	msg = mb.BuildMainMenu(chatID, currentOperations.NumberObject)
	h.sendMessage(msg)
}

func (h *CallbackHandler) getCustomers(chatID int64, currentOperation *models.Operation) {
	text := ""
	for _, customer := range currentOperation.Customers {
		text += fmt.Sprintf("№: %d\nФИО: %s\nТел.: %s\n\n", customer.UserNumber, customer.ObjCustName, customer.ObjCustPhone1)
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = menus.BackAndFinish()
	h.sendMessage(msg)
}

func (h *CallbackHandler) handleGetInfo(chatID int64) {
	op, exists := h.Operations[chatID]
	if !exists {
		return
	}

	text := fmt.Sprintf("№ объекта: %d\nНаименование: %s\nАдрес: %s",
		op.Object.AccountNumber, op.Object.Name, op.Object.Address)

	msg := tgbotapi.NewMessage(chatID, text)
	h.Bot.Send(msg)
}
