package models

import "github.com/jinleejun-corp/corelib/validators"

type DeleteFuelRefillByIDRequest struct {
	FuelRefillID int64 `validate:"required"`
}

func (req DeleteFuelRefillByIDRequest) Validate() error {
	return validators.Validate(req)
}
