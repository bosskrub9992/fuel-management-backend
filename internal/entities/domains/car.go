package domains

import "time"

type Car struct {
	ID         int64     `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (d Car) TableName() string {
	return "cars"
}
