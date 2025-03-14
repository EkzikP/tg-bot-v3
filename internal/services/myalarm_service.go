package services

import (
	"context"
	andromeda "github.com/EkzikP/sdk_andromeda_go_v2"
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
