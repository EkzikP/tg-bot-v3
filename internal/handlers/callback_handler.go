package handlers

import (
	"context"
	"fmt"
	"github.com/EkzikP/tg-bot-v3/internal/menus"
	"github.com/EkzikP/tg-bot-v3/internal/models"
	"github.com/EkzikP/tg-bot-v3/internal/services"
	"github.com/EkzikP/tg-bot-v3/internal/utils"
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
	phoneUser, _ := tgUser.Load(chatID)

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
		h.handleChecksKTS(ctx, chatID, currentOperations)
	case "MyAlarm":
		if h.haveMyAlarmRights(ctx, currentOperations, phoneUser.(string)) {
			currentOperations.Update("CurrentRequest", callback.Data)
			currentOperations.Update("CurrentMenu", "MyAlarmMenu")
			msg := menus.New().BuildMyAlarmMenu(chatID, currentOperations.NumberObject)
			h.sendMessage(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "У вас нет прав на работу с системой MyAlarm")
			h.sendMessage(msg)
			msg = menus.New().BuildMainMenu(chatID, currentOperations.NumberObject)
			h.sendMessage(msg)
		}
	case "GetUsersMyAlarm":
		currentOperation[chatID].changeValue("currentRequest", update.CallbackQuery.Data)
		msg = getUsersMyAlarm(ctx, client, confSDK, currentOperation[chatID], chatID)
		msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID

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

func (h *CallbackHandler) handleChecksKTS(ctx context.Context, chatID int64, currentOperation *models.Operation) {
	if currentOperation.CurrentRequest == "ChecksKTS" {

		resp, err := h.Andromeda.PostCheckPanic(ctx, currentOperation.Object.Id)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, err.Error())
			msg.ReplyMarkup = menus.BackAndFinish()
			h.sendMessage(msg)
			return
		}

		PostCheckPanic := map[string]string{
			"has alarm":                   "по объекту есть тревога, проверка КТС запрещена",
			"already runnig":              "по объекту уже выполняется проверка КТС",
			"success":                     "проверка КТС начата",
			"error":                       "при выполнении запроса произошла ошибка",
			"invalid checkInterval value": "для параметра checkInterval задано значение, выходящее за пределы допустимого диапазона",
		}

		if resp.Description == "already runnig" {
			text := fmt.Sprintf("По объекту уже выполняется проверка КТС.\nДождитесь автоматического завершения проверки (макс. 3 мин.) или " +
				"отправьте тревогу КТС, для завершения ранее начатой проверки.\nИ повторите попытку снова.")
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyMarkup = menus.BackAndFinish()
			h.sendMessage(msg)
			return
		} else if resp.Description != "success" {
			msg := tgbotapi.NewMessage(chatID, PostCheckPanic[resp.Description])
			msg.ReplyMarkup = menus.BackAndFinish()
			h.sendMessage(msg)
			return
		}

		currentOperation.Update("CheckPanicId", resp.CheckPanicId)

		text := fmt.Sprintf("%s\nВ течении 180 сек. нажмите кнпку КТС.\nИ нажмите кнопку \"Получить результат проверки КТС\"", PostCheckPanic[resp.Description])
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = menus.CheckKTS(resp.CheckPanicId)
		h.sendMessage(msg)
		return
	} else if currentOperation.CurrentRequest == "ResultCheckKTS" {

		resp, err := h.Andromeda.GetCheckPanic(ctx, currentOperation.CheckPanicId)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, err.Error())
			msg.ReplyMarkup = menus.CheckKTS(currentOperation.CheckPanicId)
			h.sendMessage(msg)
			return
		}

		CheckPanicResponse := map[string]string{
			"not found":                   "проверка с КТС не найдена",
			"in progress":                 "проверка КТС продолжается (не завершена): КТС не получена, тайм-аут не истек",
			"success":                     "проверка КТС успешно завершена",
			"success, interval continues": "проверка КТС успешно завершена, но продолжается отсчет интервала проверки",
			"time out":                    "проверка КТС завершена с ошибкой: истек интервал ожидания события КТС",
			"error":                       "при выполнении запроса произошла ошибка",
		}

		msg := tgbotapi.NewMessage(chatID, CheckPanicResponse[resp.Description])
		if resp.Description == "in progress" {
			msg.ReplyMarkup = menus.CheckKTS(currentOperation.CheckPanicId)
		} else {
			msg.ReplyMarkup = menus.BackAndFinish()
		}
		h.sendMessage(msg)
		return
	}
}

func (h *CallbackHandler) haveMyAlarmRights(ctx context.Context, currentOperation *models.Operation, phoneUser string) bool {

	resp, err := h.Andromeda.GetUsersMyAlarm(ctx, currentOperation.Object.Id)
	if err != nil {
		return false
	}

	currentOperation.Update("UsersMyAlarm", resp)

	var validUser bool
	for _, user := range resp {
		if user.MyAlarmPhone == phoneUser {
			validUser = true
			break
		}
	}

	if !validUser && !utils.IsEngineer(phoneUser, h.PhoneEngineer) {
		return false
	}

	return true
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
