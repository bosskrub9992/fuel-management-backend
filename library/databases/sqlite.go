package databases

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDBSqlite(fullFilePath string, gormConfig gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?parseTime=true&_loc=Asia/Bangkok",
		fullFilePath,
	)
	gormDB, err := gorm.Open(sqlite.Open(dsn), &gormConfig)
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
