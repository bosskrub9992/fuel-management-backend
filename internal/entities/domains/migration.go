package domains

import "time"

type Migration struct {
	ID        uint      `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (m Migration) TableName() string {
	return "migrations"
}
