package services

import (
	"context"

	"github.com/bosskrub9992/fuel-management-backend/internal/domains"
)

type DatabaseAdaptor interface {
	CreateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage, userIDs []int64) error
	GetCarFuelUsageWithUsers(context.Context, GetCarFuelUsageWithUsersParams) (records []FuelUsageWithUser, totalRecords int64, err error)
	GetAllUsers(context.Context) ([]domains.User, error)
	GetAllCars(context.Context) ([]domains.Car, error)
	GetLatestFuelRefill(context.Context) (*domains.FuelRefill, error)
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
