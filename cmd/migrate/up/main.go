package main

import (
	"context"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/migrations/mgpostgres"
	"github.com/bosskrub9992/fuel-management-backend/library/databases"
	"github.com/bosskrub9992/fuel-management-backend/library/zerologger"
	"github.com/rs/zerolog/log"
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
	zerologger.InitZerologExtension(cfg.Logger)

	// sort ascending
	sort.SliceStable(mgpostgres.Migrations, func(i, j int) bool {
		return mgpostgres.Migrations[i].ID < mgpostgres.Migrations[j].ID
	})

	idToMigration := make(map[uint]mgpostgres.Migration)
	for _, migration := range mgpostgres.Migrations {
		if _, found := idToMigration[migration.ID]; found {
			log.Error().Msgf("duplicate migration id: [%d]", migration.ID)
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
			log.Err(err).Send()
			return
		}
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Err(err).Send()
			}
		}()
		db, err = databases.NewGormDBPostgres(sqlDB, gorm.Config{})
		if err != nil {
			log.Err(err).Send()
			return
		}
	case "sqlite":
		db, err = databases.NewGormDBSqlite(cfg.Database.SQLite.FilePath, gorm.Config{})
		if err != nil {
			log.Err(err).Send()
			return
		}
	default:
		log.Error().Msg("invalid database type")
		return
	}

	dbMigrator := db.Migrator()
	tableMigration := domains.Migration{}.TableName()
	if !dbMigrator.HasTable(tableMigration) {
		log.Info().Msgf("not found table [%s]", tableMigration)
		if err := dbMigrator.CreateTable(&domains.Migration{}); err != nil {
			log.Err(err).Send()
			return
		}
		log.Info().Msgf("created table [%s]", tableMigration)
	}

	var inDBMigrations []domains.Migration
	if err := db.Model(&domains.Migration{}).Find(&inDBMigrations).Error; err != nil {
		log.Err(err).Send()
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
				log.Err(err).Send()
				return err
			}
			if err := migration.VerifyUp(ctx, tx); err != nil {
				log.Err(err).Send()
				return err
			}
			migrated := domains.Migration{
				ID:        migration.ID,
				CreatedAt: time.Now(),
			}
			if err := tx.Create(&migrated).Error; err != nil {
				log.Err(err).Send()
				return err
			}
			return nil
		})
		if err != nil {
			return
		}
		log.Info().Msgf("successfully migrated id: [%d] up", migration.ID)
		migratedCount++
	}

	if migratedCount == 0 {
		log.Info().Msg("no up migrations to migrate")
	} else {
		if all {
			log.Info().Msg("successfully migrated up all pending migrations")
		}
	}
}
