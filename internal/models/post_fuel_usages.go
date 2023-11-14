package models

import (
	"time"

	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type CreateFuelUsageRequest struct {
	FuelUseDate        string          `form:"fuelUseDate" validate:"max=50"`
	FuelPrice          decimal.Decimal `form:"fuelPrice"`
	KilometerBeforeUse int64           `form:"kilometerBeforeUseInput"`
	KilometerAfterUse  int64           `form:"kilometerAfterUseInput"`
	Description        string          `form:"description" validate:"max=500"`
	CarID              int64           `param:"currentCarId" validate:"required"`
	UserIDs            []int64
}

func (req CreateFuelUsageRequest) Validate() error {
	if _, err := time.Parse("2006-01-02", req.FuelUseDate); err != nil {
		return err
	}
	return validators.Validate(req)
}
