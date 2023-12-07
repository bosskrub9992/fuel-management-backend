package models

import (
	"time"

	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type PutFuelUsageRequest struct {
	FuelUsageID        int64           `json:"fuelUsageId"`
	FuelUseDate        time.Time       `json:"fuelUseDate" validate:"max=50"`
	FuelPrice          decimal.Decimal `json:"fuelPrice"`
	KilometerBeforeUse int64           `json:"kilometerBeforeUse"`
	KilometerAfterUse  int64           `json:"kilometerAfterUse"`
	Description        string          `json:"description" validate:"max=500"`
	UserIDs            []int64         `json:"userId" validate:"required"`
	IsPaid             bool            `json:"isPaid"`
	CurrentCarID       int64           `json:"currentCarId" validate:"required"`
}

func (req PutFuelUsageRequest) Validate() error {
	return validators.Validate(req)
}
