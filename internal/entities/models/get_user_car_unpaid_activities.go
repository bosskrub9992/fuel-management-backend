package models

import (
	"time"

	"github.com/bosskrub9992/fuel-management-backend/library/validators"
	"github.com/shopspring/decimal"
)

type GetUserCarUnpaidActivitiesRequest struct {
	UserID int64 `validate:"required"`
	CarID  int64 `validate:"required"`
}

func (req GetUserCarUnpaidActivitiesRequest) Validate() error {
	return validators.Validate(req)
}

type GetUserCarUnpaidActivitiesResponse struct {
	CarID       int64        `json:"carId"`
	CarName     string       `json:"carName"`
	FuelUsages  []FuelUsage  `json:"fuelUsages"`
	FuelRefills []FuelRefill `json:"fuelRefills"`
}

type FuelRefill struct {
	FuelRefillID int64           `json:"fuelRefillId"`
	RefillTime   time.Time       `json:"refillTime"`
	IsPaid       bool            `json:"isPaid"`
	TotalMoney   decimal.Decimal `json:"total_money"`
}
