package models

import (
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
	FuelUsageID     int64           `json:"fuelUsageId"`
	FuelUsageUserID int64           `json:"fuelUsageUserId"`
	FuelUseTime     string          `json:"fuelUseTime"`
	Description     string          `json:"description"`
	FuelUsers       string          `json:"fuelUsers"`
	PayEach         decimal.Decimal `json:"payEach"`
}
