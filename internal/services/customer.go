package services

import (
	"context"

	"github.com/bosskrub9992/fuel-management/internal/domains"
	"github.com/bosskrub9992/fuel-management/internal/models"

	"github.com/jinleejun-corp/corelib/slogger"
)

func (s *Service) CreateCustomer(ctx context.Context, req models.CreateCustomerRequest) (*models.CreateCustomerResponse, error) {
	customer := domains.Customer{
		ShortName: req.ShortName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	customerID, err := s.db.CreateCustomer(ctx, customer)
	if err != nil {
		slogger.Error(ctx, err.Error())
		return nil, err
	}

	return &models.CreateCustomerResponse{
		ID: customerID,
	}, nil
}
