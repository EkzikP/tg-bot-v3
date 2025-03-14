package handlers

import (
	"context"
	"fmt"
	"github.com/EkzikP/tg-bot-v3/internal/menus"
	"strconv"
	"sync"

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
	Operations map[int64]*models.Operation
}

func NewMessageHandler(bot *tgbotapi.BotAPI, cfg config.Config, store *storage.UsersStore, as *services.AndromedaService) *MessageHandler {
	return &MessageHandler{
		Bot:        bot,
		Config:     cfg,
		Storage:    store,
		Andromeda:  as,
		Operations: make(map[int64]*models.Operation),
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

		h.Operations[chatID] = models.New()
		msg := tgbotapi.NewMessage(chatID, "Введите пультовый номер объекта!")
		h.sendMessage(msg)
		return
	}

	currentOperations := h.Operations[chatID]
	if currentOperations.NumberObject == "" {
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

		if resp, ok := h.Andromeda.CheckUserRights(ctx, object, phone.(string), h.Config.PhoneEngineer); !ok {
			text := fmt.Sprintf("У вас нет прав на этот объект!\nВведите пультовый номер объекта!")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			h.sendMessage(msg)
			return
		} else {
			currentOperations.Update("NumberObject", strconv.Itoa(object.AccountNumber))
			currentOperations.Update("Object", object)
			currentOperations.Update("Customers", resp)
			currentOperations.Update("CurrentMenu", "MainMenu")
		}

		msg := tgbotapi.NewMessage(chatID, "Работа с объектом "+update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		outMsg, _ := h.Bot.Send(msg)
		pinMessage := tgbotapi.PinChatMessageConfig{
			ChatID:              chatID,
			MessageID:           outMsg.MessageID,
			DisableNotification: false,
		}
		h.Bot.Send(pinMessage)

		currentOperations.Update("CurrentMenu", "MainMenu")
		currentOperations.Update("CurrentRequest", "")

		mb := menus.New()
		msg = mb.BuildMainMenu(chatID, currentOperations.NumberObject)
		h.sendMessage(msg)
		return
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
