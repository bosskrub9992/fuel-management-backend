package main

import (
	"fmt"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/jinleejun-corp/corelib/databases"
	"github.com/jinleejun-corp/corelib/slogger"
	"gopkg.in/guregu/null.v4"
)

type Test struct {
	ID         int64     `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	CreateTime null.Time `gorm:"column:create_time"`
	UpdateTime null.Time `gorm:"column:update_time"`
}

func (d Test) TableName() string {
	return "cars"
}

func main() {
	cfg := config.New()
	logger := slogger.New(&cfg.Logger)
	db, err := databases.NewGormDBSqlite(cfg.Database.FilePath)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var cars []Test
	if err := db.Model(&Test{}).Find(&cars).Error; err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug(fmt.Sprintf("cars: %+v", cars))
}
