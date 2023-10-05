package gormadaptor

import (
	"context"

	"github.com/bosskrub9992/fuel-management/internal/domains"

	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(db *gorm.DB) *Database {
	return &Database{
		db: db,
	}
}

func (adt *Database) CreateCustomer(ctx context.Context, customer domains.Customer) (int64, error) {
	if err := adt.db.WithContext(ctx).Create(&customer).Error; err != nil {
		return 0, err
	}
	return customer.ID, nil
}
