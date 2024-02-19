package pgadaptor

import (
	"context"
	"log/slog"

	"github.com/bosskrub9992/fuel-management-backend/internal/constants"
	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
	"github.com/bosskrub9992/fuel-management-backend/library/errs"
	"gorm.io/gorm"
)

type PostgresAdaptor struct {
	db *gorm.DB
}

func NewPostgresAdaptor(db *gorm.DB) (*PostgresAdaptor, error) {
	if db == nil {
		return nil, errs.ErrNotEnoughArgForDependencyInjection
	}
	return &PostgresAdaptor{
		db: db,
	}, nil
}

func (adt *PostgresAdaptor) Transaction(ctx context.Context, fn func(ctxTx context.Context) error) error {
	return adt.db.Transaction(func(tx *gorm.DB) error {
		ctxTx := context.WithValue(ctx, constants.WithTx, tx)
		return fn(ctxTx)
	})
}

func (adt *PostgresAdaptor) dbOrTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(constants.WithTx).(*gorm.DB)
	if ok {
		return tx
	}
	return adt.db.WithContext(ctx)
}

func (adt *PostgresAdaptor) CreateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) (int64, error) {
	if err := adt.dbOrTx(ctx).Create(&fuelUsage).Error; err != nil {
		return 0, err
	}
	return fuelUsage.ID, nil
}

func (adt *PostgresAdaptor) CreateFuelUsageUsers(ctx context.Context, fuelUsageUsers []domains.FuelUsageUser) error {
	return adt.dbOrTx(ctx).
		Create(&fuelUsageUsers).
		Error
}

func (adt *PostgresAdaptor) GetFuelUsageInPagination(
	ctx context.Context,
	params services.GetFuelUsageInPaginationParams,
) (
	[]domains.FuelUsage,
	int64,
	error,
) {
	var totalCount int64
	stmt := adt.dbOrTx(ctx).
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
	err := stmt.Order("fuel_use_time DESC, id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&fuelUsages).Error
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, 0, err
	}

	return fuelUsages, totalCount, nil
}

func (adt *PostgresAdaptor) GetFuelUsageUsersByFuelUsageIDs(ctx context.Context, fuelUsageIDs []int64) ([]services.FuelUsageUser, error) {
	var fuelUsageUsers []services.FuelUsageUser
	err := adt.dbOrTx(ctx).
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

func (adt *PostgresAdaptor) GetAllUsers(ctx context.Context) ([]domains.User, error) {
	var users []domains.User
	err := adt.dbOrTx(ctx).
		Model(&domains.User{}).
		Order("nickname ASC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (adt *PostgresAdaptor) GetLatestFuelRefillByCarID(ctx context.Context, carID int64) (*domains.FuelRefill, error) {
	var fuelRefill domains.FuelRefill
	err := adt.dbOrTx(ctx).
		Model(&fuelRefill).
		Where(domains.FuelRefill{
			CarID: carID,
		}).
		Order("refill_time DESC, id DESC").
		First(&fuelRefill).Error
	if err != nil {
		return nil, err
	}
	return &fuelRefill, nil
}

func (adt *PostgresAdaptor) GetAllCars(ctx context.Context) ([]domains.Car, error) {
	var cars []domains.Car
	err := adt.dbOrTx(ctx).
		Model(&domains.Car{}).
		Find(&cars).Error
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (adt *PostgresAdaptor) GetFuelUsageByID(ctx context.Context, id int64) (*domains.FuelUsage, error) {
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

func (adt *PostgresAdaptor) GetFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) ([]services.FuelUsageUser, error) {
	var fuelUsageUsers []services.FuelUsageUser
	err := adt.dbOrTx(ctx).
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

func (adt *PostgresAdaptor) UpdateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) error {
	return adt.dbOrTx(ctx).
		Save(&fuelUsage).
		Error
}

func (adt *PostgresAdaptor) DeleteFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) error {
	return adt.dbOrTx(ctx).
		Where("fuel_usage_id = ?", fuelUsageID).
		Delete(&domains.FuelUsageUser{}).
		Error
}

func (adt *PostgresAdaptor) DeleteFuelUsageByID(ctx context.Context, id int64) error {
	return adt.dbOrTx(ctx).
		Delete(&domains.FuelUsage{}, id).
		Error
}

func (adt *PostgresAdaptor) GetFuelRefillPagination(ctx context.Context, params services.GetFuelRefillPaginationParams) ([]domains.FuelRefill, int, error) {
	stmt := adt.dbOrTx(ctx).
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
	err := stmt.Order("refill_time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&fuelRefills).Error
	if err != nil {
		return nil, 0, err
	}

	return fuelRefills, int(totalCount), nil
}

func (adt *PostgresAdaptor) CreateFuelRefill(ctx context.Context, fr domains.FuelRefill) error {
	return adt.dbOrTx(ctx).
		Create(&fr).
		Error
}

func (adt *PostgresAdaptor) GetFuelRefillByID(ctx context.Context, fuelRefillID int64) (*domains.FuelRefill, error) {
	var fr domains.FuelRefill
	err := adt.dbOrTx(ctx).
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

func (adt *PostgresAdaptor) DeleteFuelRefillByID(ctx context.Context, fuelRefillID int64) error {
	return adt.dbOrTx(ctx).
		Delete(&domains.FuelRefill{}, fuelRefillID).
		Error
}

func (adt *PostgresAdaptor) UpdateFuelRefill(ctx context.Context, fr domains.FuelRefill) error {
	return adt.dbOrTx(ctx).
		Save(&fr).
		Error
}

func (adt *PostgresAdaptor) GetLatestFuelUsageByCarID(ctx context.Context, carID int64) (*domains.FuelUsage, error) {
	var fuelUsage domains.FuelUsage
	err := adt.dbOrTx(ctx).
		Model(&fuelUsage).
		Where(domains.FuelUsage{
			CarID: carID,
		}).
		Order("fuel_use_time DESC, id DESC").
		First(&fuelUsage).Error
	if err != nil {
		return nil, err
	}
	return &fuelUsage, nil
}

func (adt *PostgresAdaptor) GetUserFuelUsagesByPaidStatus(
	ctx context.Context,
	userID int64,
	isPaid bool,
) (
	[]services.FuelUsageUserWithPayEach,
	error,
) {
	var data []services.FuelUsageUserWithPayEach
	err := adt.dbOrTx(ctx).
		Select(`fuel_usage_users.*, 
			fuel_usages.fuel_use_time, 
			fuel_usages.pay_each, 
			fuel_usages.description,
			cars.id AS car_id, 
			cars.name AS car_name`).
		Table("fuel_usages").
		Joins("INNER JOIN fuel_usage_users ON fuel_usages.id = fuel_usage_users.fuel_usage_id").
		Joins("INNER JOIN cars ON cars.id = fuel_usages.car_id").
		Where("user_id = ? AND is_paid = ?",
			userID,
			isPaid,
		).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (adt *PostgresAdaptor) UpdateUserFuelUsagePaymentStatus(ctx context.Context, userFuelUsage domains.FuelUsageUser) error {
	return adt.dbOrTx(ctx).
		Model(&domains.FuelUsageUser{}).
		Where(domains.FuelUsageUser{
			ID: userFuelUsage.ID,
		}).
		Update("is_paid", userFuelUsage.IsPaid).
		Error
}

func (adt *PostgresAdaptor) GetUserFuelUsageByUserID(ctx context.Context, userID int64) ([]domains.FuelUsageUser, error) {
	var userFuelUsages []domains.FuelUsageUser
	err := adt.dbOrTx(ctx).
		Model(&domains.FuelUsageUser{}).
		Where(domains.FuelUsageUser{
			UserID: userID,
		}).
		Find(&userFuelUsages).Error
	if err != nil {
		return nil, err
	}
	return userFuelUsages, nil
}
