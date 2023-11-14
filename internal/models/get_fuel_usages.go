package models

import (
	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type GetCarFuelUsagesRequest struct {
	CarID         int64 `param:"currentCarId" validate:"required"`
	CurrentUserID int64 `query:"currentUserId" validate:"required"`
	PageIndex     int   `query:"PageIndex"`
	PageSize      int   `query:"PageSize"`
}

func (req GetCarFuelUsagesRequest) Validate() error {
	return validators.Validate(req)
}

type GetCarFuelUsageData struct {
	TodayDate               string              `json:"todayDate"`
	LatestKilometerAfterUse int64               `json:"latestKilometerAfterUse"`
	LatestFuelPrice         decimal.Decimal     `json:"latestFuelPrice"`
	AllUsers                []User              `json:"allUsers"`
	Data                    []CarFuelUsageDatum `json:"data"`
	AllCars                 []Car               `json:"allCars"`
	TotalRecord             int64               `json:"totalRecord"`
	CurrentCar              Car                 `json:"currentCar"`
	CurrentUser             UserWithImageURL    `json:"currentUser"`
}

type CarFuelUsageDatum struct {
	ID                 int64           `json:"id"`
	FuelUseDate        string          `json:"fuelUseDate"`
	FuelPrice          decimal.Decimal `json:"fuelPrice"`
	KilometerBeforeUse int64           `json:"kilometerBeforeUse"`
	KilometerAfterUse  int64           `json:"kilometerAfterUse"`
	Description        string          `json:"description"`
	TotalMoney         decimal.Decimal `json:"totalMoney"`
	FuelUsers          string          `json:"fuelUsers"`
}

type User struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type Car struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserWithImageURL struct {
	User
	UserImageURL string `json:"userImageUrl"`
}
