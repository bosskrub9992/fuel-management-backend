package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/library/validators"
	"github.com/shopspring/decimal"
)

type CreateFuelRefillRequest struct {
	CurrentCarID          int64           `json:"currentCarId" validate:"required"`
	RefillTime            time.Time       `json:"refillTime" validate:"required"`
	KilometerBeforeRefill int64           `json:"kilometerBeforeRefill" validate:"required"`
	KilometerAfterRefill  int64           `json:"kilometerAfterRefill" validate:"required"`
	TotalMoney            decimal.Decimal `json:"totalMoney" validate:"required"`
	IsPaid                bool            `json:"isPaid"`
	RefillBy              int64           `json:"refillBy" validate:"required"`
	CurrentUserID         int64           `json:"currentUserId" validate:"required"`
}

func (req CreateFuelRefillRequest) Validate() error {
	err := validators.Validate(req)
	if req.KilometerAfterRefill <= req.KilometerBeforeRefill {
		err = errors.Join(err, fmt.Errorf("kilometerAfterRefill should > kilometerBeforeRefill"))
	}
	return err
}
