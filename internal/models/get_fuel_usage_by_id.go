package models

import (
	"time"

	"github.com/bosskrub9992/fuel-management-backend/library/validators"
	"github.com/shopspring/decimal"
)

type GetFuelUsageByIDRequest struct {
	FuelUsageID int64 `validate:"required"`
}

type GetFuelUsageByIDResponse struct {
	FuelUseTime        time.Time       `json:"fuelUseTime"`
	FuelPrice          decimal.Decimal `json:"fuelPrice"`
	FuelUsers          []GetFuelUser   `json:"fuelUsers"`
	Description        string          `json:"description"`
	KilometerBeforeUse int64           `json:"kilometerBeforeUse"`
	KilometerAfterUse  int64           `json:"kilometerAfterUse"`
	TotalMoney         decimal.Decimal `json:"totalMoney"`
	EachShouldPay      decimal.Decimal `json:"eachShouldPay"`
}

type GetFuelUser struct {
	UserID   int64  `json:"userId"`
	Nickname string `json:"nickname"`
	IsPaid   bool   `json:"isPaid"`
}

func (req GetFuelUsageByIDRequest) Validate() error {
	return validators.Validate(req)
}
