package services

import (
	"context"
	"errors"
	andromeda "github.com/EkzikP/sdk_andromeda_go_v2"
)

func (s *AndromedaService) PostCheckPanic(ctx context.Context, siteID string) (andromeda.PostCheckPanicResponse, error) {
	resp, err := s.Client.PostCheckPanic(ctx, andromeda.PostCheckPanicInput{
		SiteId: siteID,
		Config: s.Config,
	})
	if err != nil {
		err = errors.New("не удалось получить данные")
		return andromeda.PostCheckPanicResponse{}, err
	}
	return resp, nil
}

func (s *AndromedaService) GetCheckPanic(ctx context.Context, checkPanicId string) (andromeda.GetCheckPanicResponse, error) {
	resp, err := s.Client.GetCheckPanic(ctx, andromeda.GetCheckPanicInput{
		CheckPanicId: checkPanicId,
		Config:       s.Config,
	})
	if err != nil {
		return andromeda.GetCheckPanicResponse{}, err
	}
	return resp, nil
}
