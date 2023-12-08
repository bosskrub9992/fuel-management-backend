package models

import (
	"time"

	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type PutFuelRefillByIDRequest struct {
	FuelRefillID          int64           `param:"fuelRefillId" validate:"required"`
	CurrentCarID          int64           `json:"currentCarId" validate:"required"`
	RefillTime            time.Time       `json:"refillTime" validate:"required"`
	KilometerBeforeRefill int64           `json:"kilometerBeforeRefill" validate:"required"`
	KilometerAfterRefill  int64           `json:"kilometerAfterRefill" validate:"required"`
	TotalMoney            decimal.Decimal `json:"totalMoney" validate:"required"`
	IsPaid                bool            `json:"isPaid"`
	CurrentUserID         int64           `json:"currentUserId" validate:"required"`
}

func (req PutFuelRefillByIDRequest) Validate() error {
	return validators.Validate(req)
}
