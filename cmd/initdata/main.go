package main

import (
	"log/slog"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/bosskrub9992/fuel-management-backend/library/databases"
	"github.com/bosskrub9992/fuel-management-backend/library/slogger"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func main() {
	cfg := config.New()
	slog.SetDefault(slogger.New(&slogger.Config{
		IsProductionEnv: cfg.Logger.IsProductionEnv,
		MaskingFields:   cfg.Logger.MaskingFields,
		RemovingFields:  cfg.Logger.RemovingFields,
	}))
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
	db, err := databases.NewGormDBPostgres(sqlDB, gorm.Config{})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	now := time.Now()
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	cars := []domains.Car{
		{
			Name: "Mazda 2", CreateTime: now, UpdateTime: now,
		},
		{
			Name: "Ford", CreateTime: now, UpdateTime: now,
		},
	}

	users := []domains.User{
		{
			DefaultCarID:    1,
			Nickname:        "Boss",
			ProfileImageURL: "http://localhost:8080/public/BOSS.PNG",
			CreateTime:      now,
			UpdateTime:      now,
		},
		{
			DefaultCarID:    1,
			Nickname:        "Best",
			ProfileImageURL: "http://localhost:8080/public/BEST.PNG",
			CreateTime:      now,
			UpdateTime:      now,
		},
		{
			DefaultCarID:    2,
			Nickname:        "Nut",
			ProfileImageURL: "http://localhost:8080/public/NUT.PNG",
			CreateTime:      now,
			UpdateTime:      now,
		},
		{
			DefaultCarID:    1,
			Nickname:        "Pat",
			ProfileImageURL: "http://localhost:8080/public/PAT.PNG",
			CreateTime:      now,
			UpdateTime:      now,
		},
	}

	fuelRefills := []domains.FuelRefill{
		{
			CarID:                 1,
			RefillTime:            time.Date(2023, time.January, 1, 0, 0, 0, 0, loc),
			TotalMoney:            decimal.NewFromFloat(500),
			KilometerBeforeRefill: 300,
			KilometerAfterRefill:  800,
			FuelPriceCalculated:   decimal.NewFromFloat(1),
			IsPaid:                true,
			RefillBy:              1,
			CreateBy:              1,
			CreateTime:            now,
			UpdateBy:              1,
			UpdateTime:            now,
		},
		{
			CarID:                 2,
			RefillTime:            time.Date(2023, time.January, 1, 0, 0, 0, 0, loc),
			TotalMoney:            decimal.NewFromFloat(500),
			KilometerBeforeRefill: 400,
			KilometerAfterRefill:  900,
			FuelPriceCalculated:   decimal.NewFromFloat(1),
			IsPaid:                true,
			RefillBy:              1,
			CreateBy:              1,
			CreateTime:            now,
			UpdateBy:              1,
			UpdateTime:            now,
		},
	}

	fuelUsages := []domains.FuelUsage{
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 800,
			KilometerAfterUse:  700,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(100),
			PayEach:            decimal.NewFromFloat(100).DivRound(decimal.NewFromInt(2), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              2,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 900,
			KilometerAfterUse:  800,
			Description:        "mock2",
			TotalMoney:         decimal.NewFromFloat(100),
			PayEach:            decimal.NewFromFloat(100).DivRound(decimal.NewFromInt(2), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 700,
			KilometerAfterUse:  680,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			PayEach:            decimal.NewFromFloat(20).DivRound(decimal.NewFromInt(1), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 680,
			KilometerAfterUse:  660,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			PayEach:            decimal.NewFromFloat(20).DivRound(decimal.NewFromInt(1), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 660,
			KilometerAfterUse:  640,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			PayEach:            decimal.NewFromFloat(20).DivRound(decimal.NewFromInt(1), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 640,
			KilometerAfterUse:  620,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			PayEach:            decimal.NewFromFloat(20).DivRound(decimal.NewFromInt(1), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 620,
			KilometerAfterUse:  600,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			PayEach:            decimal.NewFromFloat(20).DivRound(decimal.NewFromInt(1), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 600,
			KilometerAfterUse:  580,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			PayEach:            decimal.NewFromFloat(20).DivRound(decimal.NewFromInt(1), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 580,
			KilometerAfterUse:  560,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			PayEach:            decimal.NewFromFloat(20).DivRound(decimal.NewFromInt(1), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 560,
			KilometerAfterUse:  540,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			PayEach:            decimal.NewFromFloat(20).DivRound(decimal.NewFromInt(1), 2),
			CreateTime:         now,
			UpdateTime:         now,
		},
	}

	fuelUsageUsers := []domains.FuelUsageUser{
		{
			FuelUsageID: 1, UserID: 1, IsPaid: false,
		},
		{
			FuelUsageID: 1, UserID: 2, IsPaid: true,
		},
		{
			FuelUsageID: 2, UserID: 3, IsPaid: false,
		},
		{
			FuelUsageID: 2, UserID: 4, IsPaid: true,
		},
		{
			FuelUsageID: 3, UserID: 1, IsPaid: true,
		},
		{
			FuelUsageID: 4, UserID: 1, IsPaid: true,
		},
		{
			FuelUsageID: 5, UserID: 1, IsPaid: true,
		},
		{
			FuelUsageID: 6, UserID: 1, IsPaid: true,
		},
		{
			FuelUsageID: 7, UserID: 1, IsPaid: true,
		},
		{
			FuelUsageID: 8, UserID: 1, IsPaid: true,
		},
		{
			FuelUsageID: 9, UserID: 1, IsPaid: true,
		},
		{
			FuelUsageID: 10, UserID: 1, IsPaid: true,
		},
	}

	if err := db.Create(&cars).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	if err := db.Create(&users).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	if err := db.Create(&fuelRefills).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	if err := db.Create(&fuelUsages).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	if err := db.Create(&fuelUsageUsers).Error; err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("successfully init data")
}
