package main

import (
	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/library/databases"
	"github.com/bosskrub9992/fuel-management-backend/library/zerologger"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
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
	zerologger.InitZerologExtension(cfg.Logger)
	sqlDB, err := databases.NewPostgres(&cfg.Database.Postgres)
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Err(err).Send()
		}
	}()
	db, err := databases.NewGormDBPostgres(sqlDB, gorm.Config{})
	if err != nil {
		log.Err(err).Send()
		return
	}

	var cars []Test
	if err := db.Model(&Test{}).Find(&cars).Error; err != nil {
		log.Err(err).Send()
		return
	}

	log.Debug().Msgf("cars: %+v", cars)
}
