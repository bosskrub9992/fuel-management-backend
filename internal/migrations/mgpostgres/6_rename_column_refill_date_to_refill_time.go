package mgpostgres

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         6,
		Up:         up6,
		VerifyUp:   verifyUp6,
		Down:       down6,
		VerifyDown: verifyDown6,
	})
}

func up6(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_refills RENAME COLUMN refill_date TO refill_time;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp6(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_refills": {
			ShouldHaveColumn:    {"refill_time"},
			ShouldNotHaveColumn: {"refill_date"},
		},
	}
	return validateColumnExist(migrator, validateColumnExistMap)
}

func down6(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_refills RENAME COLUMN refill_time TO refill_date;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown6(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_refills": {
			ShouldHaveColumn:    {"refill_date"},
			ShouldNotHaveColumn: {"refill_time"},
		},
	}
	return validateColumnExist(migrator, validateColumnExistMap)
}
