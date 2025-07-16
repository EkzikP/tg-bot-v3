package services

import (
	"context"
	"fmt"
	andromeda "github.com/EkzikP/sdk_andromeda_go_v2"
	"github.com/EkzikP/tg-bot-v3/internal/models"
	"github.com/EkzikP/tg-bot-v3/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (s *AndromedaService) GetUsersMyAlarm(ctx context.Context, siteID string) ([]andromeda.UserMyAlarmResponse, error) {
	resp, err := s.Client.GetUsersMyAlarm(ctx, andromeda.GetUsersMyAlarmInput{
		SiteId: siteID,
		Config: andromeda.Config{},
	})
	if err != nil {
		return []andromeda.UserMyAlarmResponse{}, err
	}
	return resp, nil
}

func (s *AndromedaService) GetUserObjectMyAlarm(ctx context.Context, phoneUser string, phoneEngineer map[string]string, update tgbotapi.Update) string {

	var text string
	if utils.IsEngineer(phoneUser, phoneEngineer) {
		if update.Message == nil {
			text = "Введите номер телефона пользователя в формате: +7xxxxxxxxxx"
			return text
		} else {
			phone, err := utils.NormalizePhone(update.Message.Text)
			if err != nil {
				text = err.Error()
				return text
			}
			phoneUser = phone
		}
	}

	resp, err := s.Client.GetUserObjectMyAlarm(ctx, andromeda.GetUserObjectMyAlarmInput{
		Phone:  phoneUser,
		Config: s.Config,
	})
	if err != nil {
		text := err.Error()
		return text
	}
	if len(resp) == 0 {
		text := "У пользователя с номером " + phoneUser + " нет объектов в приложении MyAlarm"
		return text
	}

	for _, object := range resp {
		var kts string
		var role string
		if object.IsPanic {
			kts = "Да"
		} else {
			kts = "Нет"
		}

		if object.Role == "admin" {
			role = "Администратор"
		} else {
			role = "Пользователь"
		}

		getSiteResponse, err := s.GetSite(ctx, object.ObjectGUID)
		if err != nil {
			text = err.Error()
			return text
		}

		text += fmt.Sprintf("№ объекта: %d\nНаименование: %s\nАдрес: %s\nРоль: %s\nКТС: %s\n\n", getSiteResponse.AccountNumber, getSiteResponse.Name, getSiteResponse.Address, role, kts)
	}
	return text
}

func (s *AndromedaService) PutChangeUserMyAlarm(ctx context.Context, currentOperations *models.Operation, role string) string {

	_, err := s.Client.PutChangeUserMyAlarm(ctx, andromeda.PutChangeUserMyAlarmInput{
		CustId: currentOperations.ChangedUserId,
		Role:   role,
		Config: s.Config,
	})
	if err != nil {
		var text string
		if strings.Contains(err.Error(), "User already has role,") {
			text = fmt.Sprintf("У данного пользователя уже есть права!\nДля изменения роли пользователя, необходимо \"Забрать доступ к MyAlarm\" и выдать снова с необходимыми правами.")
		} else {
			var data string
			if currentOperations.CurrentRequest == "PutDelUserMyAlarm" {
				data = "удалить пользователя из MyAlarm "
			} else {
				data = "добавить пользователя к MyAlarm"
			}
			text = fmt.Sprintf("Не удалось %s. Попробуйте позже.", data)
		}

		return text
	}

	var data string
	if currentOperations.CurrentRequest == "PutDelUserMyAlarm" {
		data = "удален"
	} else {
		data = "добавлен"
	}
	currentOperations.Update("ChangedUserId", "")
	currentOperations.Update("Role", "")
	text := "Пользователь MyAlarm успешно " + data
	return text
}

func (s *AndromedaService) PutChangeVirtualKTS(ctx context.Context, currentOperations *models.Operation, isPanic bool) error {

	err := s.Client.PutChangeKTSUserMyAlarm(ctx, andromeda.PutChangeKTSUserMyAlarmInput{
		CustId:  currentOperations.ChangedUserId,
		IsPanic: isPanic,
		Config:  s.Config,
	})
	return err
}
