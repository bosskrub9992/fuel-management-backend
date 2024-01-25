package models

import (
	"github.com/bosskrub9992/fuel-management-backend/library/validators"
	"github.com/shopspring/decimal"
)

type GetLatestFuelInfoRequest struct {
	CarID int64 `validate:"required"`
}

type GetLatestFuelInfoResponse struct {
	LatestFuelPrice         decimal.Decimal `json:"latestFuelPrice"`
	LatestKilometerAfterUse int64           `json:"latestKilometerAfterUse"`
}

func (req GetLatestFuelInfoRequest) Validate() error {
	return validators.Validate(req)
}
