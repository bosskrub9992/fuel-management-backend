package models

import "github.com/bosskrub9992/fuel-management-backend/library/validators"

type PayUserCarUnpaidActivitiesRequest struct {
	UserID           int64   `param:"userId" validate:"required"`
	CarID            int64   `param:"carId" validate:"required"`
	FuelUsageUserIDs []int64 `json:"fuelUsageUserIds"`
	FuelRefillIDs    []int64 `json:"fuelRefillIds"`
}

func (req PayUserCarUnpaidActivitiesRequest) Validate() error {
	return validators.Validate(req)
}
