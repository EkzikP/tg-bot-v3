package services

import (
	"context"
	"fmt"
	"github.com/EkzikP/sdk_andromeda_go_v2"
	"github.com/EkzikP/tg-bot-v3/internal/utils"
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

func (s *AndromedaService) CheckUserRights(ctx context.Context, object andromeda.GetSitesResponse, phone string, phoneEngineer map[string]string) ([]andromeda.GetCustomerResponse, bool) {

	resp, err := s.Client.GetCustomers(ctx, andromeda.GetCustomersInput{
		SiteId: object.Id,
		Config: s.Config,
	})
	if err != nil {
		return []andromeda.GetCustomerResponse{}, false
	}

	var useRights bool
	for _, customer := range resp {
		var phoneCustomer string
		switch len(customer.ObjCustPhone1) {
		case 12:
			phoneCustomer = customer.ObjCustPhone1
		case 11:
			phoneCustomer = "+7" + customer.ObjCustPhone1[1:]
		case 10:
			phoneCustomer = "+7" + customer.ObjCustPhone1
		default:
			phoneCustomer = ""
		}
		if phone == phoneCustomer {
			useRights = true
			break
		}
	}

	if !useRights && !utils.IsEngineer(phone, phoneEngineer) {
		return []andromeda.GetCustomerResponse{}, false
	}

	return resp, true
}
