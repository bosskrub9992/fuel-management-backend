package mgpostgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         8,
		Up:         up8,
		VerifyUp:   verifyUp8,
		Down:       down8,
		VerifyDown: verifyDown8,
	})
}

func up8(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_refills ADD COLUMN refill_by BIGINT DEFAULT NULL;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}

	var fuelRefills []domains.FuelRefill
	if err := tx.Model(&domains.FuelRefill{}).Find(&fuelRefills).Error; err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, fuelRefill := range fuelRefills {
		err := tx.Model(&domains.FuelRefill{}).
			Where(domains.FuelRefill{
				ID: fuelRefill.ID,
			}).
			Update("refill_by", fuelRefill.CreateBy).Error
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	return nil
}

func verifyUp8(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_refills": {
			ShouldHaveColumn: {"refill_by"},
		},
	}

	if err := validateColumnExist(migrator, validateColumnExistMap); err != nil {
		slog.Error(err.Error())
		return err
	}

	var fuelRefills []domains.FuelRefill
	if err := tx.Model(&domains.FuelRefill{}).Find(&fuelRefills).Error; err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, fuelRefill := range fuelRefills {
		if fuelRefill.CreateBy != fuelRefill.RefillBy {
			return fmt.Errorf("create_by should equal to refill_by, create_by:'%d', refill_by:'%d'",
				fuelRefill.CreateBy,
				fuelRefill.RefillBy,
			)
		}
	}

	return nil
}

func down8(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_refills DROP COLUMN refill_by;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown8(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_refills": {
			ShouldNotHaveColumn: {"refill_by"},
		},
	}
	return validateColumnExist(migrator, validateColumnExistMap)
}
