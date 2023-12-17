package gormadaptor

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/internal/constants"
	"github.com/bosskrub9992/fuel-management-backend/internal/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
	"github.com/shopspring/decimal"
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
	db := adt.dbOrTx(ctx)
	if err := db.WithContext(ctx).Create(&fuelUsage).Error; err != nil {
		return 0, err
	}
	return fuelUsage.ID, nil
}

func (adt *Database) CreateFuelUsageUsers(ctx context.Context, fuelUsageUsers []domains.FuelUsageUser) error {
	db := adt.dbOrTx(ctx)
	return db.WithContext(ctx).Create(&fuelUsageUsers).Error
}

type fuelUsageWithUser struct {
	ID                 int64           `gorm:"column:id"`
	FuelUseTime        time.Time       `gorm:"column:fuel_use_time"`
	FuelPrice          decimal.Decimal `gorm:"column:fuel_price"`
	KilometerBeforeUse int64           `gorm:"column:kilometer_before_use"`
	KilometerAfterUse  int64           `gorm:"column:kilometer_after_use"`
	Description        string          `gorm:"column:description"`
	TotalMoney         decimal.Decimal `gorm:"column:total_money"`
	CreateTime         time.Time       `gorm:"column:create_time"`
	UpdateTime         time.Time       `gorm:"column:update_time"`
	UserID             int64           `gorm:"column:user_id"`
	Nickname           string          `gorm:"column:nickname"`
	IsPaid             bool            `gorm:"column:is_paid"`
}

func (adt *Database) GetCarFuelUsageWithUsers(ctx context.Context, params services.GetCarFuelUsageWithUsersParams) ([]services.FuelUsageWithUser, int64, error) {
	// should filter before count
	var totalCount int64

	err := adt.dbOrTx(ctx).
		WithContext(ctx).
		Model(&domains.FuelUsage{}).
		Where(domains.FuelUsage{
			CarID: params.CarID,
		}).
		Count(&totalCount).Error
	if err != nil {
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

	var fuelUsageWithUsers []fuelUsageWithUser
	err = adt.dbOrTx(ctx).
		WithContext(ctx).
		Select("fuel_usages.*", "users.nickname", "users.id AS user_id", "fuel_usage_users.is_paid").
		Table("fuel_usages").
		Joins("INNER JOIN cars ON cars.id = fuel_usages.car_id").
		Joins("INNER JOIN fuel_usage_users ON fuel_usages.id = fuel_usage_users.fuel_usage_id").
		Joins("INNER JOIN users ON fuel_usage_users.user_id = users.id").
		Where(domains.FuelUsage{
			CarID: params.CarID,
		}).
		Order("fuel_usages.id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&fuelUsageWithUsers).Error
	if err != nil {
		return nil, 0, err
	}

	var idToFuelUsageWithUser = make(map[int64]services.FuelUsageWithUser)

	for _, f := range fuelUsageWithUsers {
		fuelUsageWithUser, found := idToFuelUsageWithUser[f.ID]
		if !found {
			idToFuelUsageWithUser[f.ID] = services.FuelUsageWithUser{
				FuelUsage: domains.FuelUsage{
					ID:                 f.ID,
					FuelUseTime:        f.FuelUseTime,
					FuelPrice:          f.FuelPrice,
					KilometerBeforeUse: f.KilometerBeforeUse,
					KilometerAfterUse:  f.KilometerAfterUse,
					Description:        f.Description,
					TotalMoney:         f.TotalMoney,
					CreateTime:         f.CreateTime,
					UpdateTime:         f.UpdateTime,
				},
				Users: []services.User{
					{
						IsPaid:   f.IsPaid,
						Nickname: f.Nickname,
					},
				},
			}
		} else {
			fuelUsageWithUser.Users = append(fuelUsageWithUser.Users, services.User{
				IsPaid:   f.IsPaid,
				Nickname: f.Nickname,
			})
			idToFuelUsageWithUser[f.ID] = fuelUsageWithUser
		}
	}

	var result []services.FuelUsageWithUser
	for _, fuelUsageWithUser := range idToFuelUsageWithUser {
		result = append(result, fuelUsageWithUser)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID > result[j].ID
	})

	return result, totalCount, nil
}

func (adt *Database) GetAllUsers(ctx context.Context) ([]domains.User, error) {
	var users []domains.User
	db := adt.dbOrTx(ctx)
	stmt := db.WithContext(ctx).Model(&domains.User{}).Order("nickname ASC")
	if err := stmt.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (adt *Database) GetLatestFuelRefill(ctx context.Context) (*domains.FuelRefill, error) {
	var fuelRefill domains.FuelRefill
	db := adt.dbOrTx(ctx)
	stmt := db.WithContext(ctx).Model(&domains.FuelRefill{})
	if err := stmt.Last(&fuelRefill).Error; err != nil {
		return nil, err
	}
	return &fuelRefill, nil
}

func (adt *Database) GetAllCars(ctx context.Context) ([]domains.Car, error) {
	var cars []domains.Car
	db := adt.dbOrTx(ctx)
	if err := db.WithContext(ctx).Model(&domains.Car{}).Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}

func (adt *Database) GetFuelUsageByID(ctx context.Context, id int64) (*domains.FuelUsage, error) {
	var fuelUsage domains.FuelUsage
	db := adt.dbOrTx(ctx)
	if err := db.Model(&domains.FuelUsage{}).First(&fuelUsage).Error; err != nil {
		return nil, err
	}
	return &fuelUsage, nil
}

func (adt *Database) GetFuelUsageUsersByFuelUsageID(ctx context.Context, fuelUsageID int64) ([]services.FuelUsageUsers, error) {
	var fuelUsageUsers []services.FuelUsageUsers
	db := adt.dbOrTx(ctx)
	stmt := db.Table("fuel_usage_users").
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
	stmt = stmt.Order("refill_time DESC").Limit(pageSize).Offset(offset)
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
		Order("fuel_use_time DESC").
		First(&fuelUsage).Error
	if err != nil {
		return nil, err
	}
	return &fuelUsage, nil
}
