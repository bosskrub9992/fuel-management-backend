package services

import (
	"context"

	"github.com/bosskrub9992/fuel-management/config"
	"github.com/bosskrub9992/fuel-management/internal/domains"
)

type DBConnector interface {
	CreateCustomer(context.Context, domains.Customer) (int64, error)
}

type Service struct {
	cfg *config.Config
	db  DBConnector
}

func New(cfg *config.Config, db DBConnector) *Service {
	return &Service{
		cfg: cfg,
		db:  db,
	}
}
