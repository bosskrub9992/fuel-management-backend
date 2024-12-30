package sqliteadaptor

import (
	"context"
	"log/slog"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
	"gorm.io/gorm"
)

type SQLiteAdaptor struct {
	db *gorm.DB
}

func NewSQLiteAdaptor(db *gorm.DB) *SQLiteAdaptor {
	return &SQLiteAdaptor{
		db: db,
	}
}

func (adt *SQLiteAdaptor) Transaction(fn func(repo *SQLiteAdaptor) error) error {
	return adt.db.Transaction(func(tx *gorm.DB) error {
		repoWithTx := NewSQLiteAdaptor(tx)
		return fn(repoWithTx)
	})
}

func (adt *SQLiteAdaptor) CreateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) (int64, error) {
	if err := adt.db.WithContext(ctx).Create(&fuelUsage).Error; err != nil {
		return 0, err
	}
	return fuelUsage.ID, nil
}

func (adt *SQLiteAdaptor) CreateFuelUsageUsers(ctx context.Context, fuelUsageUsers []domains.FuelUsageUser) error {
	return adt.db.WithContext(ctx).
		Create(&fuelUsageUsers).
		Error
}

func (adt *SQLiteAdaptor) GetFuelUsageInPagination(
	ctx context.Context,
	params services.GetFuelUsageInPaginationParams,
) (
	[]domains.FuelUsage,
	int64,
	error,
) {
	var totalCount int64
	stmt := adt.db.WithContext(ctx).
		Model(&domains.FuelUsage{}).
		Where(domains.FuelUsage{
			CarID: params.CarID,
		})

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
	err := stmt.Order("datetime(fuel_use_time) DESC, id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&fuelUsages).Error
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, 0, err
	}

	return fuelUsages, totalCount, nil
}

func (adt *SQLiteAdaptor) GetFuelUsageUsersByFuelUsageIDs(ctx context.Context, fuelUsageIDs []int64) ([]services.FuelUsageUser, error) {
	var fuelUsageUsers []services.FuelUsageUser
	err := adt.db.WithContext(ctx).
		Table("fuel_usage_users").
		Select("fuel_usage_users.*, users.nickname").
		Joins("INNER JOIN users ON users.id = fuel_usage_users.user_id").
		Where("fuel_usage_users.fuel_usage_id IN ?", fuelUsageIDs).
		Find(&fuelUsageUsers).Error
	if err != nil {
		return nil, err
	}
	return fuelUsageUsers, nil
}

func (adt *SQLiteAdaptor) GetAllUsers(ctx context.Context) ([]domains.User, error) {
	var users []domains.User
	err := adt.db.WithContext(ctx).
		Model(&domains.User{}).
		Order("nickname ASC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (adt *SQLiteAdaptor) GetLatestFuelRefillByCarID(ctx context.Context, carID int64) (*domains.FuelRefill, error) {
	var fuelRefill domains.FuelRefill
	err := adt.db.WithContext(ctx).
		Model(&fuelRefill).
		Where(domains.FuelRefill{
			CarID: carID,
		}).
		Order("datetime(refill_time) DESC, id DESC").
		First(&fuelRefill).Error
	if err != nil {
		return nil, err
	}
	return &fuelRefill, nil
}

func (adt *SQLiteAdaptor) GetAllCars(ctx context.Context) ([]domains.Car, error) {
	var cars []domains.Car
	err := adt.db.WithContext(ctx).
		Model(&domains.Car{}).
		Find(&cars).Error
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (adt *SQLiteAdaptor) GetFuelUsageByID(ctx context.Context, id int64) (*domains.FuelUsage, error) {
	var fuelUsage domains.FuelUsage
	err := adt.db.WithContext(ctx).
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

func (adt *SQLiteAdaptor) GetFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) ([]services.FuelUsageUser, error) {
	var fuelUsageUsers []services.FuelUsageUser
	err := adt.db.WithContext(ctx).
		Table("fuel_usage_users").
		Select("fuel_usage_users.*, users.nickname").
		Joins("INNER JOIN users ON users.id = fuel_usage_users.user_id").
		Where("fuel_usage_users.fuel_usage_id = ?", fuelUsageID).
		Find(&fuelUsageUsers).Error
	if err != nil {
		return nil, err
	}
	return fuelUsageUsers, nil
}

func (adt *SQLiteAdaptor) UpdateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) error {
	return adt.db.WithContext(ctx).
		Save(&fuelUsage).
		Error
}

func (adt *SQLiteAdaptor) DeleteFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) error {
	return adt.db.WithContext(ctx).
		Where("fuel_usage_id = ?", fuelUsageID).
		Delete(&domains.FuelUsageUser{}).
		Error
}

func (adt *SQLiteAdaptor) DeleteFuelUsageByID(ctx context.Context, id int64) error {
	return adt.db.WithContext(ctx).
		Delete(&domains.FuelUsage{}, id).
		Error
}

func (adt *SQLiteAdaptor) GetFuelRefillPagination(ctx context.Context, params services.GetFuelRefillPaginationParams) ([]domains.FuelRefill, int, error) {
	stmt := adt.db.WithContext(ctx).
		Model(&domains.FuelRefill{}).
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
	err := stmt.Order("datetime(refill_time) DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&fuelRefills).Error
	if err != nil {
		return nil, 0, err
	}

	return fuelRefills, int(totalCount), nil
}

func (adt *SQLiteAdaptor) CreateFuelRefill(ctx context.Context, fr domains.FuelRefill) error {
	return adt.db.WithContext(ctx).
		Create(&fr).
		Error
}

func (adt *SQLiteAdaptor) GetFuelRefillByID(ctx context.Context, fuelRefillID int64) (*domains.FuelRefill, error) {
	var fr domains.FuelRefill
	err := adt.db.WithContext(ctx).
		Model(&domains.FuelRefill{}).
		Where(domains.FuelRefill{
			ID: fuelRefillID,
		}).
		First(&fr).Error
	if err != nil {
		return nil, err
	}
	return &fr, nil
}

func (adt *SQLiteAdaptor) DeleteFuelRefillByID(ctx context.Context, fuelRefillID int64) error {
	return adt.db.WithContext(ctx).
		Delete(&domains.FuelRefill{}, fuelRefillID).
		Error
}

func (adt *SQLiteAdaptor) UpdateFuelRefill(ctx context.Context, fr domains.FuelRefill) error {
	return adt.db.WithContext(ctx).
		Save(&fr).
		Error
}

func (adt *SQLiteAdaptor) GetLatestFuelUsageByCarID(ctx context.Context, carID int64) (*domains.FuelUsage, error) {
	var fuelUsage domains.FuelUsage
	err := adt.db.WithContext(ctx).
		Model(&fuelUsage).
		Where(domains.FuelUsage{
			CarID: carID,
		}).
		Order("datetime(fuel_use_time) DESC, id DESC").
		First(&fuelUsage).Error
	if err != nil {
		return nil, err
	}
	return &fuelUsage, nil
}
