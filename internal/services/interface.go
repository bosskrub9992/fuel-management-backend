package services

import (
	"context"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
)

type DatabaseAdaptor interface {
	Transaction(ctx context.Context, fn func(ctxTx context.Context) error) error
	GetFuelUsageInPagination(ctx context.Context, params GetFuelUsageInPaginationParams) ([]domains.FuelUsage, int64, error)
	GetFuelUsageUsersByFuelUsageIDs(ctx context.Context, fuelUsageIDs []int64) ([]FuelUsageUser, error)
	GetAllUsers(context.Context) ([]domains.User, error)
	GetAllCars(context.Context) ([]domains.Car, error)
	GetLatestFuelRefillByCarID(ctx context.Context, carID int64) (*domains.FuelRefill, error)
	CreateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) (int64, error)
	CreateFuelUsageUsers(ctx context.Context, fuelUsageUsers []domains.FuelUsageUser) error
	GetFuelUsageByID(ctx context.Context, id int64) (*domains.FuelUsage, error)
	GetFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) ([]FuelUsageUser, error)
	GetLatestFuelUsageByCarID(ctx context.Context, carID int64) (*domains.FuelUsage, error)
	UpdateFuelUsage(context.Context, domains.FuelUsage) error
	DeleteFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) error
	DeleteFuelUsageByID(ctx context.Context, id int64) error
	GetFuelRefillPagination(ctx context.Context, params GetFuelRefillPaginationParams) ([]domains.FuelRefill, int, error)
	CreateFuelRefill(context.Context, domains.FuelRefill) error
	GetFuelRefillByID(ctx context.Context, fuelRefillID int64) (*domains.FuelRefill, error)
	UpdateFuelRefill(ctx context.Context, fr domains.FuelRefill) error
	DeleteFuelRefillByID(ctx context.Context, fuelRefillID int64) error
}

type FuelUsageWithUser struct {
	domains.FuelUsage
	Users []User
}

type User struct {
	IsPaid   bool
	Nickname string
}

type GetFuelUsageInPaginationParams struct {
	CarID     int64
	PageIndex int
	PageSize  int
}

type FuelUsageUser struct {
	domains.FuelUsageUser
	Nickname string `gorm:"column:nickname"`
}

type GetFuelRefillPaginationParams struct {
	CarID     int64
	PageIndex int
	PageSize  int
}
