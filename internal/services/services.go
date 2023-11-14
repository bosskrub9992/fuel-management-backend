package services

import (
	"context"

	"github.com/bosskrub9992/fuel-management/config"
	"github.com/bosskrub9992/fuel-management/internal/domains"
)

type DBConnector interface {
	CreateFuelUsage(context.Context, domains.FuelUsage) (int64, error)
	CreateFuelUsageUsers(context.Context, []domains.FuelUsageUser) error
	GetCarFuelUsageWithUsers(context.Context, GetCarFuelUsageWithUsersParams) (records []FuelUsageWithUser, totalRecords int64, err error)
	GetAllUsers(context.Context) ([]domains.User, error)
	GetAllCars(context.Context) ([]domains.Car, error)
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

type SearchField string

const (
	SearchFieldFuelUseDate = "fuelUseDate"
	SearchFieldUser        = "user"
)

type FuelUsageWithUser struct {
	domains.FuelUsage
	Users []string
}

type GetCarFuelUsageWithUsersParams struct {
	CarID     int64
	PageIndex int
	PageSize  int
}
