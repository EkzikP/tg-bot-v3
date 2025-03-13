package handlers

import (
	"context"
	"fmt"
	"github.com/EkzikP/tg-bot-v3/internal/menus"
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

func (h *MessageHandler) HandleCommand(ctx context.Context, update tgbotapi.Update, tgUsers *sync.Map) {
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
			msg.ReplyToMessageID = update.Message.MessageID
			h.sendMessage(msg)
			return
		}

		h.Operations.Store(chatID, models.New())
		msg := tgbotapi.NewMessage(chatID, "Введите пультовый номер объекта!")
		h.sendMessage(msg)
		return
	}

	currentOperations, _ := h.Operations.LoadOrStore(chatID, models.New())
	if currentOperations.(models.Operation).NumberObject == "" {
		if !h.verifyPhone(&update, tgUsers) {
			msg := h.createPhoneRequest(chatID)
			msg.ReplyToMessageID = update.Message.MessageID
			h.sendMessage(msg)
			return
		}

		if update.Message.Contact != nil {
			msg := tgbotapi.NewMessage(chatID, "Введите пультовый номер объекта!")
			msg.ReplyToMessageID = update.Message.MessageID
			h.sendMessage(msg)
			return
		}

		if message, ok := utils.CheckNumberObject(update.Message.Text); !ok {
			text := fmt.Sprintf("%s\nВведите пультовый номер объекта!", message)
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			h.sendMessage(msg)
			return
		}

		object, err := h.Andromeda.GetSite(ctx, update.Message.Text)
		if err != nil {
			text := fmt.Sprintf("%s\nВведите пультовый номер объекта!", err)
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			h.sendMessage(msg)
			return
		}

		phone, ok := tgUsers.Load(chatID)
		if !ok {
			text := fmt.Sprintf("У вас нет прав на этот объект!\nВведите пультовый номер объекта!")
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			h.sendMessage(msg)
			return
		}

		if !h.Andromeda.CheckUserRights(ctx, object.Id, phone.(string)) {
			text := fmt.Sprintf("У вас нет прав на этот объект!\nВведите пультовый номер объекта!")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			h.sendMessage(msg)
			return
		}
	}
}

func (h *MessageHandler) verifyPhone(update *tgbotapi.Update, tgUser *sync.Map) bool {
	return utils.VerifyPhone(update, tgUser, h.Storage)
}

func (h *MessageHandler) createPhoneRequest(chatID int64) tgbotapi.MessageConfig {
	menu := menus.New()
	msg := menu.RequestContact(chatID)
	return msg
}
