package handlers

import (
	"context"
	"fmt"
	"github.com/EkzikP/tg-bot-v3/internal/models"
	"github.com/EkzikP/tg-bot-v3/internal/services"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler struct {
	Bot           *tgbotapi.BotAPI
	Andromeda     *services.AndromedaService
	Operations    sync.Map
	PhoneEngineer map[string]string
}

func NewCallbackHandler(bot *tgbotapi.BotAPI, phoneEngineer map[string]string, andromeda *services.AndromedaService) *CallbackHandler {
	return &CallbackHandler{
		Bot:           bot,
		Andromeda:     andromeda,
		Operations:    sync.Map{},
		PhoneEngineer: phoneEngineer,
	}
}

func (h *CallbackHandler) HandleCallback(update tgbotapi.Update, tgUser *sync.Map) {
	callback := update.CallbackQuery
	chatID := callback.Message.Chat.ID

	switch callback.Data {
	case "GetInfoObject":
		h.handleGetInfo(chatID)
	case "ChecksKTS":
		h.handleChecksKTS(chatID)
		// Обработка других callback-ов
	}
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
