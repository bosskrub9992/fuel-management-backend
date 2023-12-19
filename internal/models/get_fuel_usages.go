package models

import (
	"github.com/jinleejun-corp/corelib/validators"
	"github.com/shopspring/decimal"
)

type GetFuelUsagesRequest struct {
	CurrentCarID  int64 `query:"currentCarId" validate:"required"`
	CurrentUserID int64 `query:"currentUserId" validate:"required"`
	PageIndex     int   `query:"pageIndex"`
	PageSize      int   `query:"pageSize"`
}

type GetFuelUsagesResponse struct {
	FuelUsageData []FuelUsageDatum `json:"fuelUsageData"`
	TotalRecord   int64            `json:"totalRecord"`
	TotalPage     int64            `json:"totalPage"`
}

type FuelUsageDatum struct {
	ID                 int64           `json:"id"`
	FuelUseTime        string          `json:"fuelUseTime"`
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

func (req GetFuelUsagesRequest) Validate() error {
	return validators.Validate(req)
}
