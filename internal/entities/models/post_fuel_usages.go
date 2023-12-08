package models

import (
	"errors"
	"time"

	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type CreateFuelUsageRequest struct {
	CurrentCarID       int64           `json:"currentCarId" validate:"required"`
	FuelUseTime        time.Time       `json:"fuelUseTime"`
	FuelPrice          decimal.Decimal `json:"fuelPrice"`
	FuelUsers          []FuelUser      `json:"fuelUsers" validate:"min=1"`
	Description        string          `json:"description" validate:"max=500"`
	KilometerBeforeUse int64           `json:"kilometerBeforeUse"`
	KilometerAfterUse  int64           `json:"kilometerAfterUse"`
}

type FuelUser struct {
	UserID int64 `json:"userId" validate:"required"`
	IsPaid bool  `json:"isPaid" validate:"required"`
}

func (req CreateFuelUsageRequest) Validate() error {
	err := validators.Validate(req)
	if req.KilometerBeforeUse < req.KilometerAfterUse {
		err = errors.Join(err, errors.New("kmBeforeUse should be > kmAfterUse"))
	}

	var hasDuplicateFuelUser bool
	numberOfFuelUser := len(req.FuelUsers)
	if numberOfFuelUser > 0 {
		uniqueUserID := make(map[int64]bool)
		for _, fuelUser := range req.FuelUsers {
			uniqueUserID[fuelUser.UserID] = true
		}
		if numberOfFuelUser > len(uniqueUserID) {
			hasDuplicateFuelUser = true
		}
	}

	if hasDuplicateFuelUser {
		err = errors.Join(err, errors.New("fuelUsers has duplicate userId"))
	}

	return err
}
