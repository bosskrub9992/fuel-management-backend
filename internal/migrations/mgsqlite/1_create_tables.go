package mgsqlite

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

func init() {
	Migrations = append(Migrations, Migration{
		ID:         1,
		Up:         up1,
		VerifyUp:   verifyUp1,
		Down:       down1,
		VerifyDown: verifyDown1,
	})
}

func up1(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS fuel_usages (
			id INTEGER PRIMARY KEY,
			fuel_use_time DATETIME NOT NULL,
			fuel_price DECIMAL(10,3) NOT NULL,
			kilometer_before_use INT NOT NULL,
			kilometer_after_use INT NOT NULL,
			description VARCHAR(500),
			total_money DECIMAL(10,3) NOT NULL,
			create_time DATETIME NOT NULL,
			update_time DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			nickname VARCHAR(500) NOT NULL,
			create_time DATETIME NOT NULL,
			update_time DATETIME NOT null
		);`,
		`CREATE TABLE IF NOT EXISTS fuel_usage_users (
			id INTEGER PRIMARY KEY,
			fuel_usage_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS fuel_refills (
			id INTEGER PRIMARY KEY,
			refill_date DATETIME NOT NULL,
			total_money DECIMAL(10,3) NOT NULL,
			kilometer_before_refill INT NOT NULL,
			kilometer_after_refill INT NOT NULL,
			fuel_price_calculated DECIMAL(10,3) NOT NULL,
			refill_by VARCHAR(50) NOT NULL,
			create_time DATETIME NOT NULL,
			update_time DATETIME NOT NULL
		);`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyUp1(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	return tableShouldExist(migrator,
		"fuel_usages",
		"users",
		"fuel_usage_users",
		"fuel_refills",
	)
}

func down1(ctx context.Context, tx *gorm.DB) error {
	sqlStatements := []string{
		`DROP TABLE IF EXISTS fuel_usages;`,
		`DROP TABLE IF EXISTS users;`,
		`DROP TABLE IF EXISTS fuel_usage_users;`,
		`DROP TABLE IF EXISTS fuel_refills;`,
	}
	for index, sqlStatement := range sqlStatements {
		if err := tx.WithContext(ctx).Exec(sqlStatement).Error; err != nil {
			slog.Error(err.Error(), "index", index)
			return err
		}
	}
	return nil
}

func verifyDown1(ctx context.Context, tx *gorm.DB) error {
	migrator := tx.Migrator()
	return tableShouldNotExist(migrator,
		"fuel_usages",
		"users",
		"fuel_usage_users",
		"fuel_refills",
	)
}
