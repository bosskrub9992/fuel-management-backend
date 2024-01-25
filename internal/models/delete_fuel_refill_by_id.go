package models

import "github.com/bosskrub9992/fuel-management-backend/library/validators"

type DeleteFuelRefillByIDRequest struct {
	FuelRefillID int64 `validate:"required"`
}

func (req DeleteFuelRefillByIDRequest) Validate() error {
	return validators.Validate(req)
}
