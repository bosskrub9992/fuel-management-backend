package domains

import "time"

type User struct {
	ID              int64     `gorm:"column:id"`
	DefaultCarID    int64     `gorm:"column:default_car_id"`
	Nickname        string    `gorm:"column:nickname"`
	ProfileImageURL string    `gorm:"column:profile_image_url"`
	CreateTime      time.Time `gorm:"column:create_time"`
	UpdateTime      time.Time `gorm:"column:update_time"`
}

func (d User) TableName() string {
	return "users"
}
