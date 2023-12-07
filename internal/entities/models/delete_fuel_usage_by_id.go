package models

import (
	"github.com/jinleejun-corp/corelib/validators"
)

type DeleteFuelUsageByIDRequest struct {
	FuelUsageID int64 `validate:"required"`
}

func (req DeleteFuelUsageByIDRequest) Validate() error {
	return validators.Validate(req)
}
