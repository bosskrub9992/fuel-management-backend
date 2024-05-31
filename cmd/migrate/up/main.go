package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/migrations/mgpostgres"
	"github.com/bosskrub9992/fuel-management-backend/library/databases"
	"github.com/bosskrub9992/fuel-management-backend/library/slogger"
	"gorm.io/gorm"
)

func main() {
	args := os.Args

	var all bool
	if len(args) > 1 && args[1] == "all" {
		all = true
	}

	cfg := config.New()
	ctx := context.Background()
	slog.SetDefault(slogger.New(&slogger.Config{
		IsProductionEnv: cfg.Logger.IsProductionEnv,
		MaskingFields:   cfg.Logger.MaskingFields,
		RemovingFields:  cfg.Logger.RemovingFields,
	}))

	// sort ascending
	sort.SliceStable(mgpostgres.Migrations, func(i, j int) bool {
		return mgpostgres.Migrations[i].ID < mgpostgres.Migrations[j].ID
	})

	idToMigration := make(map[uint]mgpostgres.Migration)
	for _, migration := range mgpostgres.Migrations {
		if _, found := idToMigration[migration.ID]; found {
			slog.Error(fmt.Sprintf("duplicate migration id: [%d]",
				migration.ID,
			))
			return
		}
		idToMigration[migration.ID] = migration
	}

	var db *gorm.DB
	var err error

	switch strings.ToLower(cfg.Database.Use) {
	case "postgres":
		sqlDB, err := databases.NewPostgres(&cfg.Database.Postgres)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		defer func() {
			if err := sqlDB.Close(); err != nil {
				slog.Error(err.Error())
			}
		}()
		db, err = databases.NewGormDBPostgres(sqlDB, gorm.Config{})
		if err != nil {
			slog.Error(err.Error())
			return
		}
	case "sqlite":
		db, err = databases.NewGormDBSqlite(cfg.Database.SQLite.FilePath, gorm.Config{})
		if err != nil {
			slog.Error(err.Error())
			return
		}
	default:
		slog.Error("invalid database type")
		return
	}

	dbMigrator := db.Migrator()
	tableMigration := domains.Migration{}.TableName()
	if !dbMigrator.HasTable(tableMigration) {
		slog.Info(fmt.Sprintf("not found table [%s]", tableMigration))
		if err := dbMigrator.CreateTable(&domains.Migration{}); err != nil {
			slog.Error(err.Error())
			return
		}
		slog.Info(fmt.Sprintf("created table [%s]", tableMigration))
	}

	var inDBMigrations []domains.Migration
	if err := db.Model(&domains.Migration{}).Find(&inDBMigrations).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	idToInDBMigration := make(map[uint]domains.Migration)
	for _, migration := range inDBMigrations {
		idToInDBMigration[migration.ID] = migration
	}

	var migratedCount int

	for _, migration := range mgpostgres.Migrations {
		if !all && migratedCount == 1 {
			break
		}
		if _, found := idToInDBMigration[migration.ID]; found {
			continue
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Up(ctx, tx); err != nil {
				slog.Error(err.Error())
				return err
			}
			if err := migration.VerifyUp(ctx, tx); err != nil {
				slog.Error(err.Error())
				return err
			}
			migrated := domains.Migration{
				ID:        migration.ID,
				CreatedAt: time.Now(),
			}
			if err := tx.Create(&migrated).Error; err != nil {
				slog.Error(err.Error())
				return err
			}
			return nil
		})
		if err != nil {
			return
		}
		slog.Info(fmt.Sprintf("successfully migrated id: [%d] up", migration.ID))
		migratedCount++
	}

	if migratedCount == 0 {
		slog.Info("no up migrations to migrate")
	} else {
		if all {
			slog.Info("successfully migrated up all pending migrations")
		}
	}
}
