package models

import (
	"time"

	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type GetFuelRefillRequest struct {
	CurrentCarID int64 `query:"currentCarId" validate:"required"`
	PageIndex    int   `query:"pageIndex"`
	PageSize     int   `query:"pageSize"`
}

func (req GetFuelRefillRequest) Validate() error {
	return validators.Validate(req)
}

type GetFuelRefillResponse struct {
	Data        []FuelRefillDatum `json:"data"`
	TotalRecord int               `json:"totalRecord"`
}

type FuelRefillDatum struct {
	RefillTime            time.Time       `json:"refillTime"`
	KilometerBeforeRefill int64           `json:"kilometerBeforeRefill"`
	KilometerAfterRefill  int64           `json:"kilometerAfterRefill"`
	TotalMoney            decimal.Decimal `json:"totalMoney"`
	FuelPriceCalculated   decimal.Decimal `json:"fuelPriceCalculated"`
	IsPaid                bool            `json:"isPaid"`
}
