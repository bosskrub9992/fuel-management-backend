package mgpostgres

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         4,
		Up:         up4,
		VerifyUp:   verifyUp4,
		Down:       down4,
		VerifyDown: verifyDown4,
	})
}

func up4(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_refills ADD COLUMN is_paid BOOL;`,
		`ALTER TABLE fuel_usages ADD COLUMN is_paid BOOL;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp4(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	tableToColumns := map[string][]string{
		"fuel_refills": {"is_paid"},
		"fuel_usages":  {"is_paid"},
	}
	return tableShouldHaveColumns(migrator, tableToColumns)
}

func down4(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`ALTER TABLE fuel_refills DROP COLUMN is_paid;`,
		`ALTER TABLE fuel_usages DROP COLUMN is_paid;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown4(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	tableToColumns := map[string][]string{
		"fuel_refills": {"is_paid"},
		"fuel_usages":  {"is_paid"},
	}
	return tableShouldNotHaveColumns(migrator, tableToColumns)
}
