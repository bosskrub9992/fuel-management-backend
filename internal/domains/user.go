package domains

import "time"

type User struct {
	ID         int64     `gorm:"column:id"`
	Nickname   string    `gorm:"column:nickname"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (d User) TableName() string {
	return "users"
}
