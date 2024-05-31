package mgpostgres

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         5,
		Up:         up5,
		VerifyUp:   verifyUp5,
		Down:       down5,
		VerifyDown: verifyDown5,
	})
}

func up5(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_usages DROP COLUMN is_paid;`,
		`ALTER TABLE fuel_usage_users ADD COLUMN is_paid BOOL;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp5(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_usage_users": {
			ShouldHaveColumn: {"is_paid"},
		},
		"fuel_usages": {
			ShouldNotHaveColumn: {"is_paid"},
		},
	}
	return validateColumnExist(migrator, validateColumnExistMap)
}

func down5(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_usages ADD COLUMN is_paid BOOL;`,
		`ALTER TABLE fuel_usage_users DROP COLUMN is_paid;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown5(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_usages": {
			ShouldHaveColumn: {"is_paid"},
		},
		"fuel_usage_users": {
			ShouldNotHaveColumn: {"is_paid"},
		},
	}
	return validateColumnExist(migrator, validateColumnExistMap)
}
