package models

import (
	"time"

	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type PutFuelUsageRequest struct {
	FuelUsageID        int64           `param:"fuelUsageId" validate:"required"`
	CurrentCarID       int64           `json:"currentCarId" validate:"required"`
	FuelUseTime        time.Time       `json:"fuelUseTime"`
	FuelPrice          decimal.Decimal `json:"fuelPrice"`
	FuelUsers          []FuelUser      `json:"fuelUsers" validate:"min=1"`
	Description        string          `json:"description" validate:"max=500"`
	KilometerBeforeUse int64           `json:"kilometerBeforeUse"`
	KilometerAfterUse  int64           `json:"kilometerAfterUse"`
}

func (req PutFuelUsageRequest) Validate() error {
	return validators.Validate(req)
}
