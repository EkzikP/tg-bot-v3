package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/EkzikP/sdk_andromeda_go_v2"
)

type AndromedaService struct {
	Client *andromeda.Client
	Config andromeda.Config
}

func New(cfg andromeda.Config) *AndromedaService {
	return &AndromedaService{
		Client: andromeda.NewClient(),
		Config: cfg,
	}
}

func (s *AndromedaService) GetSite(ctx context.Context, siteID string) (andromeda.GetSitesResponse, error) {
	resp, err := s.Client.GetSites(ctx, andromeda.GetSitesInput{
		Id:     siteID,
		Config: s.Config,
	})
	if err != nil {
		return andromeda.GetSitesResponse{}, fmt.Errorf("failed to get site: %w", err)
	}
	return resp, nil
}

func (s *AndromedaService) CheckUserRights(ctx context.Context, siteID string, phone string) (bool, error) {
	// Реализация проверки прав
}
