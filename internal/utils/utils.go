package utils

import (
	"fmt"
	andromeda "github.com/EkzikP/sdk_andromeda_go_v2"
	"github.com/EkzikP/tg-bot-v3/internal/models"
	"strconv"
)

func CheckNumberObject(text string) (string, bool) {

	num, err := strconv.Atoi(text)
	if err != nil {
		return "Номер объекта введен некорректно!", false
	}

	if num < 1 || num > 9999 {
		return "Номер объекта введен некорректно!", false
	}
	return "", true
}

func ListUsersMyAlarm(currentOperations *models.Operation) string {

	text := ""
	if len(currentOperations.UsersMyAlarm) == 0 {
		text = "Не найдено ни одного пользователя"
	}
	for _, user := range currentOperations.UsersMyAlarm {
		var kts string
		var role string
		if user.IsPanic {
			kts = "Да"
		} else {
			kts = "Нет"
		}

		if user.Role == "admin" {
			role = "Администратор"
		} else {
			role = "Пользователь"
		}

		var customer andromeda.GetCustomerResponse
		for _, customer = range currentOperations.Customers {
			if customer.Id == user.CustomerID {
				break
			}
		}

		text += fmt.Sprintf("ФИО: %s\nТел.: %s\nРоль: %s\nКТС: %s\n\n", customer.ObjCustName, user.MyAlarmPhone, role, kts)
	}
	return text
}
