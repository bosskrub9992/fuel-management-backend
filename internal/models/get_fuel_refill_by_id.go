package models

import (
	"time"

	"github.com/bosskrub9992/fuel-management-backend/library/validators"
	"github.com/shopspring/decimal"
)

type GetFuelRefillByIDRequest struct {
	FuelRefillID int64 `validate:"required"`
}

type GetFuelRefillByIDResponse struct {
	ID                    int64           `json:"id"`
	RefillTime            time.Time       `json:"refillTime"`
	KilometerBeforeRefill int64           `json:"kilometerBeforeRefill"`
	KilometerAfterRefill  int64           `json:"kilometerAfterRefill"`
	TotalMoney            decimal.Decimal `json:"totalMoney"`
	FuelPriceCalculated   decimal.Decimal `json:"fuelPriceCalculated"`
	IsPaid                bool            `json:"isPaid"`
}

func (req GetFuelRefillByIDRequest) Validate() error {
	return validators.Validate(req)
}
