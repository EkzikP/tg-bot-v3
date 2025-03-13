package handlers

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/EkzikP/sdk_andromeda_go_v2"
	"github.com/EkzikP/tg-bot-v3/internal/config"
	"github.com/EkzikP/tg-bot-v3/internal/models"
	"github.com/EkzikP/tg-bot-v3/internal/services"
	"github.com/EkzikP/tg-bot-v3/internal/storage"
	"github.com/EkzikP/tg-bot-v3/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandler struct {
	Bot        *tgbotapi.BotAPI
	Config     config.Config
	Storage    *storage.UsersStore
	Andromeda  *services.AndromedaService
	Operations sync.Map
}

func NewMessageHandler(bot *tgbotapi.BotAPI, cfg config.Config, store *storage.UsersStore, as *services.AndromedaService) *MessageHandler {
	return &MessageHandler{
		Bot:        bot,
		Config:     cfg,
		Storage:    store,
		Andromeda:  as,
		Operations: sync.Map{},
	}
}

func (h *MessageHandler) sendMessage(msg tgbotapi.MessageConfig) {
	h.Bot.Send(msg)
}

func (h *MessageHandler) HandleCommand(update tgbotapi.Update, tgUsers *sync.Map) {
	chatID := update.Message.Chat.ID

	if !update.Message.Chat.IsPrivate() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Бот работает только в приватных чатах")
		msg.ReplyToMessageID = update.Message.MessageID
		h.sendMessage(msg)
		return
	}

	if update.Message.IsCommand() {

		if !h.verifyPhone(&update, tgUsers) {
			msg := h.createPhoneRequest(chatID)
			h.sendMessage(msg)
			return
		}

		h.Operations[chatID] = models.New()
		msg := tgbotapi.NewMessage(chatID, "Введите пультовый номер объекта!")
		h.sendMessage(msg)
		return
	}
}

func (h *MessageHandler) verifyPhone(update *tgbotapi.Update, tgUser *sync.Map) bool {
	return utils.VerifyPhone(update, tgUser, h.Storage)
}

func (h *MessageHandler) createPhoneRequest(chatID int64) tgbotapi.MessageConfig {
	// Создание запроса телефона
}
