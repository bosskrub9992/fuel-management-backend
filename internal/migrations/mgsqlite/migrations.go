package mgsqlite

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Migration struct {
	ID         uint
	Up         func(ctx context.Context, tx *gorm.DB) error
	VerifyUp   func(ctx context.Context, tx *gorm.DB) error
	Down       func(ctx context.Context, tx *gorm.DB) error
	VerifyDown func(ctx context.Context, tx *gorm.DB) error
}

var Migrations = []Migration{}

func columnShouldExist(migrator gorm.Migrator, table string, columns ...string) (err error) {
	for _, column := range columns {
		if !migrator.HasColumn(table, column) {
			err = errors.Join(err, fmt.Errorf("column %q should exists in table %q",
				column,
				table,
			))
		}
	}
	return err
}

func columnShouldNotExist(migrator gorm.Migrator, table string, columns ...string) (err error) {
	for _, column := range columns {
		if migrator.HasColumn(table, column) {
			err = errors.Join(err, fmt.Errorf("column %q should not exists in table %q",
				column,
				table,
			))
		}
	}
	return err
}

func tableShouldExist(migrator gorm.Migrator, tables ...string) (err error) {
	for _, table := range tables {
		if !migrator.HasTable(table) {
			err = errors.Join(err, fmt.Errorf("table %q should exists",
				table,
			))
		}
	}
	return err
}

func tableShouldNotExist(migrator gorm.Migrator, tables ...string) (err error) {
	for _, table := range tables {
		if migrator.HasTable(table) {
			err = errors.Join(err, fmt.Errorf("table %q should not exists",
				table,
			))
		}
	}
	return err
}

func tableShouldHaveColumns(migrator gorm.Migrator, tableToColumns map[string][]string) (err error) {
	for table, columns := range tableToColumns {
		for _, column := range columns {
			if !migrator.HasColumn(table, column) {
				err = errors.Join(err, fmt.Errorf("table %q should have column %q",
					table,
					column,
				))
			}
		}
	}
	return err
}

func tableShouldNotHaveColumns(migrator gorm.Migrator, tableToColumns map[string][]string) (err error) {
	for table, columns := range tableToColumns {
		for _, column := range columns {
			if migrator.HasColumn(table, column) {
				err = errors.Join(err, fmt.Errorf("table %q should not have column %q",
					table,
					column,
				))
			}
		}
	}
	return err
}

type ColumnType bool

const (
	ShouldHaveColumn    ColumnType = true
	ShouldNotHaveColumn ColumnType = false
)

func validateColumnExist(migrator gorm.Migrator, validateColumnExistMap map[string]map[ColumnType][]string) (err error) {
	for table, columnTypeToColumns := range validateColumnExistMap {
		columns, found := columnTypeToColumns[ShouldHaveColumn]
		if found {
			err = errors.Join(err,
				columnShouldExist(migrator, table, columns...),
			)
		}
		columns, found = columnTypeToColumns[ShouldNotHaveColumn]
		if found {
			err = errors.Join(err,
				columnShouldNotExist(migrator, table, columns...),
			)
		}
	}
	return err
}
