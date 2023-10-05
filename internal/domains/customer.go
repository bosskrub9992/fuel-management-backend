package domains

import "gopkg.in/guregu/null.v4"

type Customer struct {
	ID        int64       `gorm:"column:id"`
	ShortName string      `gorm:"column:short_name"`
	FirstName null.String `gorm:"column:first_name"`
	LastName  null.String `gorm:"column:last_name"`
}

func (c Customer) TableName() string {
	return "customer"
}
