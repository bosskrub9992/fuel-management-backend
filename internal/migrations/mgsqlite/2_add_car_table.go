package mgsqlite

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         2,
		Up:         up2,
		VerifyUp:   verifyUp2,
		Down:       down2,
		VerifyDown: verifyDown2,
	})
}

func up2(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS cars (
			id INTEGER PRIMARY KEY,
			name VARCHAR(500) NOT NULL,
			create_time DATETIME NOT NULL,
			update_time DATETIME NOT NULL
		);`,

		`ALTER TABLE fuel_refills ADD COLUMN car_id BIGINT NOT NULL;`, 
		`ALTER TABLE fuel_refills ADD COLUMN update_by BIGINT NOT NULL;`, 
		`ALTER TABLE fuel_refills ADD COLUMN create_by BIGINT NOT NULL;`, 

		`ALTER TABLE fuel_refills DROP COLUMN refill_by;`,

		`ALTER TABLE fuel_usages ADD COLUMN car_id BIGINT NOT NULL;`,

		`ALTER TABLE users ADD COLUMN default_car_id BIGINT NOT NULL;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp2(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	validateColumnExistMap := map[string]map[ColumnType][]string{
		"cars": {
			ShouldHaveColumn: {"id", "name", "create_time", "update_time"},
		},
		"fuel_refills": {
			ShouldHaveColumn:    {"car_id", "update_by", "create_by"},
			ShouldNotHaveColumn: {"refill_by"},
		},
		"fuel_usages": {
			ShouldHaveColumn: {"car_id"},
		},
		"users": {
			ShouldHaveColumn: {"default_car_id"},
		},
	}
	return validateColumnExist(migrator, validateColumnExistMap)
}

func down2(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`DROP TABLE IF EXISTS cars;`,

		`ALTER TABLE fuel_refills DROP COLUMN car_id`,
		`ALTER TABLE fuel_refills DROP COLUMN update_by;`,
		`ALTER TABLE fuel_refills DROP COLUMN create_by;`,

		`ALTER TABLE fuel_refills ADD COLUMN refill_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM';`,

		`ALTER TABLE fuel_usages DROP COLUMN car_id;`,
		`ALTER TABLE users DROP COLUMN default_car_id;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown2(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()

	if err := tableShouldNotExist(migrator, "cars"); err != nil {
		slog.Error(err.Error())
		return err
	}

	validateColumnExistMap := map[string]map[ColumnType][]string{
		"fuel_refills": {
			ShouldHaveColumn:    {"refill_by"},
			ShouldNotHaveColumn: {"car_id", "update_by", "create_by"},
		},
		"fuel_usages": {
			ShouldNotHaveColumn: {"car_id"},
		},
		"users": {
			ShouldNotHaveColumn: {"default_car_id"},
		},
	}
	return validateColumnExist(migrator, validateColumnExistMap)
}
