package gormadaptor

import (
	"context"
	"log/slog"

	"github.com/bosskrub9992/fuel-management-backend/internal/constants"
	"github.com/bosskrub9992/fuel-management-backend/internal/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
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

func (adt *Database) Transaction(ctx context.Context, fn func(ctxTx context.Context) error) error {
	return adt.db.Transaction(func(tx *gorm.DB) error {
		ctxTx := context.WithValue(ctx, constants.WithTx, tx)
		return fn(ctxTx)
	})
}

func (adt *Database) dbOrTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(constants.WithTx).(*gorm.DB)
	if ok {
		return tx
	}
	return adt.db
}

func (adt *Database) CreateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) (int64, error) {
	if err := adt.dbOrTx(ctx).WithContext(ctx).Create(&fuelUsage).Error; err != nil {
		return 0, err
	}
	return fuelUsage.ID, nil
}

func (adt *Database) CreateFuelUsageUsers(ctx context.Context, fuelUsageUsers []domains.FuelUsageUser) error {
	return adt.dbOrTx(ctx).
		WithContext(ctx).
		Create(&fuelUsageUsers).
		Error
}

func (adt *Database) GetFuelUsageInPagination(
	ctx context.Context,
	params services.GetFuelUsageInPaginationParams,
) (
	[]domains.FuelUsage,
	int64,
	error,
) {
	var totalCount int64
	stmt := adt.dbOrTx(ctx).
		WithContext(ctx).
		Model(&domains.FuelUsage{}).
		Where(domains.FuelUsage{
			CarID: params.CarID,
		})

	// should filter here if have any

	if err := stmt.Count(&totalCount).Error; err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, 0, err
	}

	pageIndex := params.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 0
	}
	offset := (pageIndex - 1) * pageSize

	var fuelUsages []domains.FuelUsage
	stmt = stmt.Debug().Order("datetime(fuel_use_time) DESC, id DESC").
		Limit(pageSize).
		Offset(offset)
	if err := stmt.Find(&fuelUsages).Error; err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, 0, err
	}

	return fuelUsages, totalCount, nil
}

func (adt *Database) GetFuelUsageUsersByFuelUsageIDs(ctx context.Context, fuelUsageIDs []int64) ([]services.FuelUsageUser, error) {
	var fuelUsageUsers []services.FuelUsageUser
	stmt := adt.dbOrTx(ctx).
		WithContext(ctx).
		Table("fuel_usage_users").
		Select("fuel_usage_users.*, users.nickname").
		Joins("INNER JOIN users ON users.id = fuel_usage_users.user_id").
		Where("fuel_usage_users.fuel_usage_id IN ?", fuelUsageIDs)
	if err := stmt.Find(&fuelUsageUsers).Error; err != nil {
		return nil, err
	}
	return fuelUsageUsers, nil
}

func (adt *Database) GetAllUsers(ctx context.Context) ([]domains.User, error) {
	var users []domains.User
	stmt := adt.dbOrTx(ctx).
		WithContext(ctx).
		Model(&domains.User{}).
		Order("nickname ASC")
	if err := stmt.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (adt *Database) GetLatestFuelRefill(ctx context.Context) (*domains.FuelRefill, error) {
	var fuelRefill domains.FuelRefill
	stmt := adt.dbOrTx(ctx).
		WithContext(ctx).
		Model(&fuelRefill)
	if err := stmt.Last(&fuelRefill).Error; err != nil {
		return nil, err
	}
	return &fuelRefill, nil
}

func (adt *Database) GetAllCars(ctx context.Context) ([]domains.Car, error) {
	var cars []domains.Car
	err := adt.dbOrTx(ctx).
		WithContext(ctx).
		Model(&domains.Car{}).
		Find(&cars).Error
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (adt *Database) GetFuelUsageByID(ctx context.Context, id int64) (*domains.FuelUsage, error) {
	var fuelUsage domains.FuelUsage
	err := adt.dbOrTx(ctx).
		Model(&fuelUsage).
		Where(domains.FuelUsage{
			ID: id,
		}).
		First(&fuelUsage).Error
	if err != nil {
		return nil, err
	}
	return &fuelUsage, nil
}

func (adt *Database) GetFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) ([]services.FuelUsageUser, error) {
	var fuelUsageUsers []services.FuelUsageUser
	stmt := adt.dbOrTx(ctx).
		Table("fuel_usage_users").
		Select("fuel_usage_users.*, users.nickname").
		Joins("INNER JOIN users ON users.id = fuel_usage_users.user_id").
		Where("fuel_usage_users.fuel_usage_id = ?", fuelUsageID)
	if err := stmt.Find(&fuelUsageUsers).Error; err != nil {
		return nil, err
	}
	return fuelUsageUsers, nil
}

func (adt *Database) UpdateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) error {
	db := adt.dbOrTx(ctx)
	return db.Save(&fuelUsage).Error
}

func (adt *Database) DeleteFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) error {
	db := adt.dbOrTx(ctx)
	stmt := db.Where("fuel_usage_id = ?", fuelUsageID)
	return stmt.Delete(&domains.FuelUsageUser{}).Error
}

func (adt *Database) DeleteFuelUsageByID(ctx context.Context, id int64) error {
	db := adt.dbOrTx(ctx)
	return db.Delete(&domains.FuelUsage{}, id).Error
}

func (adt *Database) GetFuelRefillPagination(ctx context.Context, params services.GetFuelRefillPaginationParams) ([]domains.FuelRefill, int, error) {
	db := adt.dbOrTx(ctx)
	stmt := db.Model(&domains.FuelRefill{}).
		Where("car_id = ?", params.CarID)

	var totalCount int64
	if err := stmt.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	pageIndex := params.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 0
	}
	offset := (pageIndex - 1) * pageSize

	var fuelRefills []domains.FuelRefill
	stmt = stmt.Order("datetime(refill_time) DESC").Limit(pageSize).Offset(offset)
	if err := stmt.Find(&fuelRefills).Error; err != nil {
		return nil, 0, err
	}

	return fuelRefills, int(totalCount), nil
}

func (adt *Database) CreateFuelRefill(ctx context.Context, fr domains.FuelRefill) error {
	db := adt.dbOrTx(ctx)
	return db.Create(&fr).Error
}

func (adt *Database) GetFuelRefillByID(ctx context.Context, fuelRefillID int64) (*domains.FuelRefill, error) {
	var fr domains.FuelRefill
	db := adt.dbOrTx(ctx)
	stmt := db.Model(&domains.FuelRefill{}).
		Where(domains.FuelRefill{
			ID: fuelRefillID,
		})
	if err := stmt.First(&fr).Error; err != nil {
		return nil, err
	}
	return &fr, nil
}

func (adt *Database) DeleteFuelRefillByID(ctx context.Context, fuelRefillID int64) error {
	db := adt.dbOrTx(ctx)
	return db.Delete(&domains.FuelRefill{}, fuelRefillID).Error
}

func (adt *Database) UpdateFuelRefill(ctx context.Context, fr domains.FuelRefill) error {
	db := adt.dbOrTx(ctx)
	return db.Save(&fr).Error
}

func (adt *Database) GetLatestFuelUsage(ctx context.Context) (*domains.FuelUsage, error) {
	var fuelUsage domains.FuelUsage
	err := adt.dbOrTx(ctx).
		Model(&fuelUsage).
		Order("datetime(fuel_use_time) DESC, id DESC").
		First(&fuelUsage).Error
	if err != nil {
		return nil, err
	}
	return &fuelUsage, nil
}
