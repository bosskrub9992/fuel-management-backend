package services

import (
	"context"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
)

type DatabaseAdaptor interface {
	Transaction(ctx context.Context, fn func(ctxTx context.Context) error) error
	GetCarFuelUsageWithUsers(context.Context, GetCarFuelUsageWithUsersParams) (records []FuelUsageWithUser, totalRecords int64, err error)
	GetAllUsers(context.Context) ([]domains.User, error)
	GetAllCars(context.Context) ([]domains.Car, error)
	GetLatestFuelRefill(context.Context) (*domains.FuelRefill, error)
	CreateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) (int64, error)
	CreateFuelUsageUsers(ctx context.Context, fuelUsageUsers []domains.FuelUsageUser) error
	GetFuelUsageByID(ctx context.Context, id int64) (*domains.FuelUsage, error)
	GetFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) ([]FuelUsageUsers, error)
	UpdateFuelUsage(context.Context, domains.FuelUsage) error
	DeleteFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) error
	DeleteFuelUsageByID(ctx context.Context, id int64) error
}

type FuelUsageWithUser struct {
	domains.FuelUsage
	Users []string
}

type GetCarFuelUsageWithUsersParams struct {
	CarID     int64
	PageIndex int
	PageSize  int
}

type FuelUsageUsers struct {
	domains.FuelUsageUser
	Nickname string `gorm:"column:nickname"`
}
