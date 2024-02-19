package models

import "github.com/bosskrub9992/fuel-management-backend/library/validators"

type BulkUpdateUserFuelUsagePaymentStatusRequest struct {
	UserID         int64 `validate:"required"`
	UserFuelUsages []struct {
		ID     int64 `json:"id"`
		IsPaid bool  `json:"isPaid"`
	} `json:"userFuelUsages"`
}

func (req BulkUpdateUserFuelUsagePaymentStatusRequest) Validate() error {
	return validators.Validate(req)
}
