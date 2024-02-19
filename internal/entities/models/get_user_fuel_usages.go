package models

import (
	"time"

	"github.com/bosskrub9992/fuel-management-backend/library/validators"
	"github.com/shopspring/decimal"
)

type GetUserFuelUsagesRequest struct {
	UserID int64 `validate:"required"`
	IsPaid bool
}

func (req GetUserFuelUsagesRequest) Validate() error {
	return validators.Validate(req)
}

type GetUserFuelUsagesResponse struct {
	UserFuelUsages []UserFuelUsage `json:"userFuelUsages"`
}

type UserFuelUsage struct {
	Car        CarInfo     `json:"car"`
	FuelUsages []FuelUsage `json:"fuelUsages"`
}

type CarInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FuelUsage struct {
	ID              int64           `json:"id"`
	FuelUsageUserID int64           `json:"fuelUsageUserId"`
	FuelUseTime     time.Time       `json:"fuelUseTime"`
	PayEach         decimal.Decimal `json:"payEach"`
}
