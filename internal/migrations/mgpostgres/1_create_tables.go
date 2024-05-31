package mgpostgres

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
			id SERIAL PRIMARY KEY NOT NULL,
			fuel_use_time TIMESTAMP WITH TIME ZONE NOT NULL,
			fuel_price DECIMAL(10,3) NOT NULL,
			kilometer_before_use INT NOT NULL,
			kilometer_after_use INT NOT NULL,
			description VARCHAR(500),
			total_money DECIMAL(10,3) NOT NULL,
			create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			update_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY NOT NULL,
			nickname VARCHAR(500) NOT NULL,
			create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			update_time TIMESTAMP WITH TIME ZONE NOT null DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS fuel_usage_users (
			id SERIAL PRIMARY KEY NOT NULL,
			fuel_usage_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS fuel_refills (
			id SERIAL PRIMARY KEY NOT NULL,
			refill_date TIMESTAMP WITH TIME ZONE NOT NULL,
			total_money DECIMAL(10,3) NOT NULL,
			kilometer_before_refill INT NOT NULL,
			kilometer_after_refill INT NOT NULL,
			fuel_price_calculated DECIMAL(10,3) NOT NULL,
			refill_by VARCHAR(50) NOT NULL,
			create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			update_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
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
