package models

import (
	"time"

	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type CreateFuelUsageRequest struct {
	FuelUseDate        string          `form:"fuelUseDate"`
	FuelPrice          decimal.Decimal `form:"fuelPrice"`
	KilometerBeforeUse int64           `form:"kilometerBeforeUse"`
	KilometerAfterUse  int64           `form:"kilometerAfterUse"`
	Description        string          `form:"description" validate:"max=500"`
	UserIDs            []int64
}

func (req CreateFuelUsageRequest) Validate() error {
	if _, err := time.Parse("2006-01-02", req.FuelUseDate); err != nil {
		return err
	}
	return validators.Validate(req)
}

type GetAllFuelUsageData struct {
	TodayDate               string
	LatestKilometerAfterUse int64
	LatestFuelPrice         decimal.Decimal
	AllUsers                []User
	Data                    []GetAllFuelUsageDatum
}

type GetAllFuelUsageDatum struct {
	ID                 int64
	FuelUseDate        string
	FuelPrice          decimal.Decimal
	KilometerBeforeUse int64
	KilometerAfterUse  int64
	Description        string
	TotalMoney         decimal.Decimal
	FuelUsers          string
}

type User struct {
	ID       int64
	Nickname string
}
