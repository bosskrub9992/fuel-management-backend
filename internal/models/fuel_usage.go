package models

type FuelUsageRecord struct {
	ID              int64
	CreateTime      string
	StartFuelAmount float64
	EndFuelAmount   float64
}
