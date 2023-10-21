package gormadaptor

import (
	"context"
	"sort"
	"time"

	"github.com/bosskrub9992/fuel-management/internal/domains"
	"github.com/bosskrub9992/fuel-management/internal/services"
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

func (adt *Database) CreateFuelUsage(ctx context.Context, fuelUsage domains.FuelUsage) (int64, error) {
	if err := adt.db.WithContext(ctx).Create(&fuelUsage).Error; err != nil {
		return 0, err
	}
	return fuelUsage.ID, nil
}

func (adt *Database) CreateFuelUsageUsers(ctx context.Context, fuelUsageUsers []domains.FuelUsageUser) error {
	if err := adt.db.WithContext(ctx).Create(&fuelUsageUsers).Error; err != nil {
		return err
	}
	return nil
}

type fuelUsageWithUser struct {
	ID                 int64           `gorm:"column:id"`
	FuelUseDate        time.Time       `gorm:"column:fuel_use_date"`
	FuelPrice          decimal.Decimal `gorm:"column:fuel_price"`
	KilometerBeforeUse int64           `gorm:"column:kilometer_before_use"`
	KilometerAfterUse  int64           `gorm:"column:kilometer_after_use"`
	Description        string          `gorm:"column:description"`
	TotalMoney         decimal.Decimal `gorm:"column:total_money"`
	CreateTime         time.Time       `gorm:"column:create_time"`
	UpdateTime         time.Time       `gorm:"column:update_time"`
	UserID             int64           `gorm:"column:user_id"`
	Nickname           string          `gorm:"column:nickname"`
}

func (adt *Database) GetAllFuelUsageWithUsers(ctx context.Context) ([]services.FuelUsageWithUser, error) {
	var fuelUsageWithUsers []fuelUsageWithUser
	err := adt.db.WithContext(ctx).
		Select("fuel_usages.*", "users.nickname", "users.id AS user_id").
		Table("fuel_usages").
		Joins("INNER JOIN fuel_usage_users ON fuel_usages.id = fuel_usage_users.fuel_usage_id").
		Joins("INNER JOIN users ON fuel_usage_users.user_id = users.id").
		Order("fuel_usages.id DESC").
		Find(&fuelUsageWithUsers).Error
	if err != nil {
		return nil, err
	}

	var idToFuelUsageWithUsers = make(map[int64]services.FuelUsageWithUser)
	for _, f := range fuelUsageWithUsers {
		if _, found := idToFuelUsageWithUsers[f.ID]; !found {
			idToFuelUsageWithUsers[f.ID] = services.FuelUsageWithUser{
				FuelUsage: domains.FuelUsage{
					ID:                 f.ID,
					FuelUseDate:        f.FuelUseDate,
					FuelPrice:          f.FuelPrice,
					KilometerBeforeUse: f.KilometerBeforeUse,
					KilometerAfterUse:  f.KilometerAfterUse,
					Description:        f.Description,
					TotalMoney:         f.TotalMoney,
					CreateTime:         f.CreateTime,
					UpdateTime:         f.UpdateTime,
				},
				Users: []string{},
			}
		}
		temp := idToFuelUsageWithUsers[f.ID]
		temp.Users = append(temp.Users, f.Nickname)
		idToFuelUsageWithUsers[f.ID] = temp
	}

	var result []services.FuelUsageWithUser
	for _, fuelUsageWithUser := range idToFuelUsageWithUsers {
		result = append(result, fuelUsageWithUser)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID > result[j].ID
	})

	return result, nil
}

func (adt *Database) GetAllUsers(ctx context.Context) ([]domains.User, error) {
	var users []domains.User
	if err := adt.db.WithContext(ctx).Model(&domains.User{}).Order("nickname ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (adt *Database) GetLatestFuelRefill(ctx context.Context) (*domains.FuelRefill, error) {
	var fuelRefill domains.FuelRefill
	if err := adt.db.WithContext(ctx).Model(&domains.FuelRefill{}).Last(&fuelRefill).Error; err != nil {
		return nil, err
	}
	return &fuelRefill, nil
}
