package domains

import (
	"time"

	"github.com/shopspring/decimal"
)

type FuelRefill struct {
	ID                    int64           `gorm:"column:id"`
	CarID                 int64           `gorm:"column:car_id"`
	RefillTime            time.Time       `gorm:"column:refill_time"`
	TotalMoney            decimal.Decimal `gorm:"column:total_money"`
	KilometerBeforeRefill int64           `gorm:"column:kilometer_before_refill"`
	KilometerAfterRefill  int64           `gorm:"column:kilometer_after_refill"`
	FuelPriceCalculated   decimal.Decimal `gorm:"column:fuel_price_calculated"`
	IsPaid                bool            `gorm:"column:is_paid"`
	RefillBy              int64           `gorm:"column:refill_by"`
	CreateBy              int64           `gorm:"column:create_by"`
	CreateTime            time.Time       `gorm:"column:create_time"`
	UpdateBy              int64           `gorm:"column:update_by"`
	UpdateTime            time.Time       `gorm:"column:update_time"`
}

func (d FuelRefill) TableName() string {
	return "fuel_refills"
}
