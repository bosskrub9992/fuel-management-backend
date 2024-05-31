package mgpostgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         7,
		Up:         up7,
		VerifyUp:   verifyUp7,
		Down:       down7,
		VerifyDown: verifyDown7,
	})
}

func up7(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_usages ADD COLUMN pay_each DECIMAL(10,3) DEFAULT NULL;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}

	type fuelUsageWithPassengerCount struct {
		ID             int64           `gorm:"column:id"`
		TotalMoney     decimal.Decimal `gorm:"column:total_money"`
		PassengerCount int64           `gorm:"column:passenger_count"`
	}

	var fuelUsageWithPassengerCounts []fuelUsageWithPassengerCount
	err := tx.Raw(
		`SELECT fuel_usages.id, fuel_usages.total_money, COUNT(fuel_usage_users.id) AS passenger_count
		FROM fuel_usages 
		INNER JOIN fuel_usage_users ON fuel_usages.id = fuel_usage_users.fuel_usage_id
		GROUP BY fuel_usages.id`,
	).Find(&fuelUsageWithPassengerCounts).Error
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, f := range fuelUsageWithPassengerCounts {
		payEach := f.TotalMoney.DivRound(decimal.NewFromInt(f.PassengerCount), 2)

		err := tx.Model(&domains.FuelUsage{}).
			Where(domains.FuelUsage{
				ID: f.ID,
			}).
			Update("pay_each", payEach).Error
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	return nil
}

func verifyUp7(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_usages": {
			ShouldHaveColumn: {"pay_each"},
		},
	}

	if err := validateColumnExist(migrator, validateColumnExistMap); err != nil {
		slog.Error(err.Error())
		return err
	}

	var countNullPayEachRow int64
	err := tx.Model(&domains.FuelUsage{}).
		Where("pay_each IS NULL").
		Count(&countNullPayEachRow).Error
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if countNullPayEachRow != 0 {
		return fmt.Errorf("all fuel_usages records should have value in 'pay_each' field")
	}

	return nil
}

func down7(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_usages DROP COLUMN pay_each;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown7(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_usages": {
			ShouldNotHaveColumn: {"pay_each"},
		},
	}
	return validateColumnExist(migrator, validateColumnExistMap)
}
