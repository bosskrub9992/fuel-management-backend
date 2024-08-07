package models

import (
	"github.com/bosskrub9992/fuel-management-backend/library/validators"
	"github.com/shopspring/decimal"
)

type GetFuelRefillRequest struct {
	CurrentCarID int64 `query:"currentCarId" validate:"required"`
	PageIndex    int   `query:"pageIndex"`
	PageSize     int   `query:"pageSize"`
}

type GetFuelRefillResponse struct {
	FuelRefillData []FuelRefillDatum `json:"fuelRefillData"`
	TotalRecord    int               `json:"totalRecord"`
	TotalPage      int               `json:"totalPage"`
}

type FuelRefillDatum struct {
	ID                    int64           `json:"id"`
	RefillTime            string          `json:"refillTime"`
	KilometerBeforeRefill int64           `json:"kilometerBeforeRefill"`
	KilometerAfterRefill  int64           `json:"kilometerAfterRefill"`
	TotalMoney            decimal.Decimal `json:"totalMoney"`
	FuelPriceCalculated   decimal.Decimal `json:"fuelPriceCalculated"`
	IsPaid                bool            `json:"isPaid"`
	RefillBy              int64           `json:"refillBy"`
}

func (req GetFuelRefillRequest) Validate() error {
	return validators.Validate(req)
}
