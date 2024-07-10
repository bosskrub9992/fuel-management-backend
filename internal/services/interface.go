package services

import (
	"context"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/shopspring/decimal"
)

type DatabaseAdaptor interface {
	Transaction(ctx context.Context, fn func(ctxTx context.Context) error) error
	GetFuelUsageInPagination(ctx context.Context, params GetFuelUsageInPaginationParams) ([]domains.FuelUsage, int64, error)
	GetUserFuelUsagesByPaidStatus(ctx context.Context, userID int64, isPaid bool, carID int64) ([]FuelUsageUserWithPayEach, error)
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
	UpdateUserFuelUsagePaymentStatus(ctx context.Context, userFuelUsage domains.FuelUsageUser) error
	GetUserFuelUsageByUserID(ctx context.Context, userID int64) ([]domains.FuelUsageUser, error)
	IsUserOwnAllFuelUsageUser(ctx context.Context, userID int64, fuelUsageUserIds []int64) (bool, error)
	DeleteFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) error
	DeleteFuelUsageByID(ctx context.Context, id int64) error
	GetFuelRefillPagination(ctx context.Context, params GetFuelRefillPaginationParams) ([]domains.FuelRefill, int, error)
	CreateFuelRefill(context.Context, domains.FuelRefill) error
	GetFuelRefillByID(ctx context.Context, fuelRefillID int64) (*domains.FuelRefill, error)
	IsUserOwnAllFuelRefills(ctx context.Context, userID int64, fuelRefillIDs []int64) (bool, error)
	GetUserUnpaidFuelRefills(ctx context.Context, userID int64, carID int64) ([]domains.FuelRefill, error)
	UpdateFuelRefill(ctx context.Context, fr domains.FuelRefill) error
	DeleteFuelRefillByID(ctx context.Context, fuelRefillID int64) error
	PayFuelRefills(ctx context.Context, fuelRefillIDs []int64) error
	PayFuelUsageUsers(ctx context.Context, fuelUsageUserIds []int64) error
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

type FuelUsageUserWithPayEach struct {
	domains.FuelUsageUser
	PayEach     decimal.Decimal `gorm:"column:pay_each"`
	FuelUseTime time.Time       `gorm:"column:fuel_use_time"`
	Description string          `gorm:"column:description"`
	CarID       int64           `gorm:"column:car_id"`
	CarName     string          `gorm:"column:car_name"`
}
