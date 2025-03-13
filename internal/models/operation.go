package models

import (
	"github.com/EkzikP/sdk_andromeda_go_v2"
)

type Operation struct {
	NumberObject   string
	Object         andromeda.GetSitesResponse
	Customers      []andromeda.GetCustomerResponse
	UsersMyAlarm   []andromeda.UserMyAlarmResponse
	CurrentRequest string
	CurrentMenu    string
	CheckPanicId   string
	ChangedUserId  string
	Role           string
}

func New() *Operation {
	return &Operation{}
}

func (o *Operation) Update(field string, value interface{}) {
	switch field {
	case "NumberObject":
		o.NumberObject = value.(string)
	case "Object":
		o.Object = value.(andromeda.GetSitesResponse)
	case "Customers":
		o.Customers = value.([]andromeda.GetCustomerResponse)
	case "UsersMyAlarm":
		o.UsersMyAlarm = value.([]andromeda.UserMyAlarmResponse)
	case "CurrentRequest":
		o.CurrentRequest = value.(string)
	case "CurrentMenu":
		o.CurrentMenu = value.(string)
	case "CheckPanicId":
		o.CheckPanicId = value.(string)
	case "ChangedUserId":
		o.ChangedUserId = value.(string)
	case "Role":
		o.Role = value.(string)
	}
}
