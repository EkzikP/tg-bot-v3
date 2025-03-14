package models

import (
	"github.com/EkzikP/sdk_andromeda_go_v2"
	"sync"
)

type Operation struct {
	mu             sync.Mutex
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
		o.mu.Lock()
		o.NumberObject = value.(string)
		o.mu.Unlock()
	case "Object":
		o.mu.Lock()
		o.Object = value.(andromeda.GetSitesResponse)
		o.mu.Unlock()
	case "Customers":
		o.mu.Lock()
		o.Customers = value.([]andromeda.GetCustomerResponse)
		o.mu.Unlock()
	case "UsersMyAlarm":
		o.mu.Lock()
		o.UsersMyAlarm = value.([]andromeda.UserMyAlarmResponse)
		o.mu.Unlock()
	case "CurrentRequest":
		o.mu.Lock()
		o.CurrentRequest = value.(string)
		o.mu.Unlock()
	case "CurrentMenu":
		o.mu.Lock()
		o.CurrentMenu = value.(string)
		o.mu.Unlock()
	case "CheckPanicId":
		o.mu.Lock()
		o.CheckPanicId = value.(string)
		o.mu.Unlock()
	case "ChangedUserId":
		o.mu.Lock()
		o.ChangedUserId = value.(string)
		o.mu.Unlock()
	case "Role":
		o.mu.Lock()
		o.Role = value.(string)
		o.mu.Unlock()
	}
}
