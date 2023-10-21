package domains

type FuelUsageUser struct {
	ID          int64 `gorm:"column:id"`
	FuelUsageID int64 `gorm:"column:fuel_usage_id"`
	UserID      int64 `gorm:"column:user_id"`
}

func (d FuelUsageUser) TableName() string {
	return "fuel_usage_users"
}
