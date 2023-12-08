package models

import (
	"time"

	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type GetFuelRefillByIDRequest struct {
	FuelRefillID int64 `validate:"required"`
}

func (req GetFuelRefillByIDRequest) Validate() error {
	return validators.Validate(req)
}

type GetFuelRefillByIDResponse struct {
	RefillTime            time.Time       `json:"refillTime"`
	KilometerBeforeRefill int64           `json:"kilometerBeforeRefill"`
	KilometerAfterRefill  int64           `json:"kilometerAfterRefill"`
	TotalMoney            decimal.Decimal `json:"totalMoney"`
	FuelPriceCalculated   decimal.Decimal `json:"fuelPriceCalculated"`
	IsPaid                bool            `json:"isPaid"`
}
