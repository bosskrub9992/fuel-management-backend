package models

import (
	"github.com/bosskrub9992/fuel-management-backend/library/validators"
)

type DeleteFuelUsageByIDRequest struct {
	FuelUsageID int64 `validate:"required"`
}

func (req DeleteFuelUsageByIDRequest) Validate() error {
	return validators.Validate(req)
}
