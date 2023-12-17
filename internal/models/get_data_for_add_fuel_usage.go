package models

import "github.com/shopspring/decimal"

type GetLatestFuelInfoResponse struct {
	LatestFuelPrice         decimal.Decimal `json:"latestFuelPrice"`
	LatestKilometerAfterUse int64           `json:"latestKilometerAfterUse"`
}
