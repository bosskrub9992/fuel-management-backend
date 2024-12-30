package pgadaptor

import (
	"context"
	"log/slog"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
	"gorm.io/gorm"
)

type PostgresAdaptor struct {
	db *gorm.DB
}

func NewPostgresAdaptor(db *gorm.DB) *PostgresAdaptor {
	return &PostgresAdaptor{
		db: db,
	}
}

func (adt *PostgresAdaptor) Transaction(fn func(repo services.DatabaseAdaptor) error) error {
	return adt.db.Transaction(func(tx *gorm.DB) error {
		repoWithTx := NewPostgresAdaptor(tx)
		return fn(repoWithTx)
	})
}

func (adt *PostgresAdaptor) CreateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) (int64, error) {
	if err := adt.db.WithContext(ctx).Create(&fuelUsage).Error; err != nil {
		return 0, err
	}
	return fuelUsage.ID, nil
}

func (adt *PostgresAdaptor) CreateFuelUsageUsers(ctx context.Context, fuelUsageUsers []domains.FuelUsageUser) error {
	return adt.db.WithContext(ctx).
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

func (adt *PostgresAdaptor) GetAllUsers(ctx context.Context) ([]domains.User, error) {
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

func (adt *PostgresAdaptor) GetLatestFuelRefillByCarID(ctx context.Context, carID int64) (*domains.FuelRefill, error) {
	var fuelRefill domains.FuelRefill
	err := adt.db.WithContext(ctx).
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
	err := adt.db.WithContext(ctx).
		Model(&domains.Car{}).
		Find(&cars).Error
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (adt *PostgresAdaptor) GetFuelUsageByID(ctx context.Context, id int64) (*domains.FuelUsage, error) {
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

func (adt *PostgresAdaptor) GetFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) ([]services.FuelUsageUser, error) {
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

func (adt *PostgresAdaptor) UpdateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) error {
	return adt.db.WithContext(ctx).
		Save(&fuelUsage).
		Error
}

func (adt *PostgresAdaptor) DeleteFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) error {
	return adt.db.WithContext(ctx).
		Where("fuel_usage_id = ?", fuelUsageID).
		Delete(&domains.FuelUsageUser{}).
		Error
}

func (adt *PostgresAdaptor) DeleteFuelUsageByID(ctx context.Context, id int64) error {
	return adt.db.WithContext(ctx).
		Delete(&domains.FuelUsage{}, id).
		Error
}

func (adt *PostgresAdaptor) GetFuelRefillPagination(ctx context.Context, params services.GetFuelRefillPaginationParams) ([]domains.FuelRefill, int, error) {
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
	return adt.db.WithContext(ctx).
		Create(&fr).
		Error
}

func (adt *PostgresAdaptor) GetFuelRefillByID(ctx context.Context, fuelRefillID int64) (*domains.FuelRefill, error) {
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

func (adt *PostgresAdaptor) DeleteFuelRefillByID(ctx context.Context, fuelRefillID int64) error {
	return adt.db.WithContext(ctx).
		Delete(&domains.FuelRefill{}, fuelRefillID).
		Error
}

func (adt *PostgresAdaptor) UpdateFuelRefill(ctx context.Context, fr domains.FuelRefill) error {
	return adt.db.WithContext(ctx).
		Save(&fr).
		Error
}

func (adt *PostgresAdaptor) GetLatestFuelUsageByCarID(ctx context.Context, carID int64) (*domains.FuelUsage, error) {
	var fuelUsage domains.FuelUsage
	err := adt.db.WithContext(ctx).
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
	carID int64,
) (
	[]services.FuelUsageUserWithPayEach,
	error,
) {
	var data []services.FuelUsageUserWithPayEach
	q := adt.db.WithContext(ctx).
		Select(`fuu.*,
			fu.fuel_use_time,
			fu.description,
			fu.pay_each,
			cars.id AS car_id,
			cars.name AS car_name`).
		Table("fuel_usages AS fu").
		Joins("INNER JOIN fuel_usage_users AS fuu ON fu.id = fuu.fuel_usage_id").
		Joins("INNER JOIN cars ON cars.id = fu.car_id").
		Where("fuu.user_id = ? AND fuu.is_paid = ?",
			userID,
			isPaid,
		)

	if carID != 0 {
		q = q.Where("cars.id = ?", carID)
	}

	if err := q.Order("fu.fuel_use_time, fu.id ASC").Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (adt *PostgresAdaptor) UpdateUserFuelUsagePaymentStatus(ctx context.Context, userFuelUsage domains.FuelUsageUser) error {
	return adt.db.WithContext(ctx).
		Model(&domains.FuelUsageUser{}).
		Where(domains.FuelUsageUser{
			ID: userFuelUsage.ID,
		}).
		Update("is_paid", userFuelUsage.IsPaid).
		Error
}

func (adt *PostgresAdaptor) GetUserFuelUsageByUserID(ctx context.Context, userID int64) ([]domains.FuelUsageUser, error) {
	var userFuelUsages []domains.FuelUsageUser
	err := adt.db.WithContext(ctx).
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

func (adt *PostgresAdaptor) GetUserUnpaidFuelRefills(ctx context.Context, userID int64, carID int64) ([]domains.FuelRefill, error) {
	var unpaidFuelRefills []domains.FuelRefill
	err := adt.db.WithContext(ctx).
		Model(&domains.FuelRefill{}).
		Where("fuel_refills.is_paid = false").
		Where("fuel_refills.refill_by = ?", userID).
		Where("fuel_refills.car_id = ?", carID).
		Order("refill_time ASC").
		Find(&unpaidFuelRefills).Error
	if err != nil {
		return nil, err
	}
	return unpaidFuelRefills, nil
}

func (adt *PostgresAdaptor) PayFuelRefills(ctx context.Context, fuelRefillIDs []int64) error {
	err := adt.db.WithContext(ctx).
		Model(&domains.FuelRefill{}).
		Where("id IN ?", fuelRefillIDs).
		Update("is_paid", true).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (adt *PostgresAdaptor) PayFuelUsageUsers(ctx context.Context, fuelUsageUserIds []int64) error {
	err := adt.db.WithContext(ctx).
		Model(&domains.FuelUsageUser{}).
		Where("id IN ?", fuelUsageUserIds).
		Update("is_paid", true).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (adt *PostgresAdaptor) IsUserOwnAllFuelRefills(ctx context.Context, userID int64, fuelRefillIDs []int64) (bool, error) {
	var count int64
	err := adt.db.WithContext(ctx).
		Model(&domains.FuelRefill{}).
		Where("id IN ?", fuelRefillIDs).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if int(count) == len(fuelRefillIDs) {
		return true, nil
	}
	return false, nil
}

func (adt *PostgresAdaptor) IsUserOwnAllFuelUsageUser(ctx context.Context, userID int64, fuelUsageUserIds []int64) (bool, error) {
	var count int64
	err := adt.db.WithContext(ctx).
		Model(&domains.FuelUsageUser{}).
		Where("id IN ?", fuelUsageUserIds).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if int(count) == len(fuelUsageUserIds) {
		return true, nil
	}
	return false, nil
}
