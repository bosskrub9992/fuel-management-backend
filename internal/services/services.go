package services

import (
	"context"

	"github.com/bosskrub9992/fuel-management/config"
	"github.com/bosskrub9992/fuel-management/internal/domains"
)

type DBConnector interface {
	CreateFuelUsage(context.Context, domains.FuelUsage) (int64, error)
	CreateFuelUsageUsers(context.Context, []domains.FuelUsageUser) error
	GetAllFuelUsageWithUsers(context.Context) ([]FuelUsageWithUser, error)
	GetAllUsers(context.Context) ([]domains.User, error)
	GetLatestFuelRefill(context.Context) (*domains.FuelRefill, error)
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

type FuelUsageWithUser struct {
	domains.FuelUsage
	Users []string
}
