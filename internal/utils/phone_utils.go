package utils

import (
	"errors"
	"github.com/EkzikP/tg-bot-v3/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"sync"
)

func NormalizePhone(phone string) (string, error) {
	clean := strings.TrimPrefix(phone, "+")

	switch len(clean) {
	case 11:
		return "+7" + clean[1:], nil
	case 10:
		return "+7" + clean, nil
	default:
		return "", errors.New("invalid phone format")
	}
}

func IsEngineer(phone string, engineers map[string]string) bool {
	_, exists := engineers[phone]
	return exists
}

func VerifyPhone(update *tgbotapi.Update, tgUsers *sync.Map, store *storage.UsersStore) bool {

	chatID := update.Message.Chat.ID

	if update.Message.Contact == nil {
		if _, ok := tgUsers.Load(chatID); !ok {
			phone, err := store.Get(chatID)
			if err != nil {
				return false
			}
			tgUsers.Store(chatID, phone)
			return true
		}
		return true
	}

	contactPhone, err := NormalizePhone(update.Message.Contact.PhoneNumber)
	if err != nil {
		return false
	}

	if phone, ok := tgUsers.Load(chatID); !ok {
		tgUsers.Store(chatID, contactPhone)
		err = store.Add(chatID, contactPhone)
		if err != nil {
			return false
		}
		return true
	} else if phone != contactPhone {
		tgUsers.Store(chatID, contactPhone)
		err = store.Add(chatID, contactPhone)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
