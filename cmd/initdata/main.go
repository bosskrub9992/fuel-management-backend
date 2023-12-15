package main

import (
	"log/slog"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/internal/domains"
	"github.com/jinleejun-corp/corelib/databases"
	"github.com/jinleejun-corp/corelib/slogger"
	"github.com/shopspring/decimal"
)

func main() {
	cfg := config.New()
	logger := slogger.New(&cfg.Logger)
	db, err := databases.NewGormDBSqlite(cfg.Database.FilePath)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	now := time.Now()
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	cars := []domains.Car{
		{
			ID: 1, Name: "Mazda 2", CreateTime: now, UpdateTime: now,
		},
		{
			ID: 2, Name: "Ford", CreateTime: now, UpdateTime: now,
		},
	}

	users := []domains.User{
		{
			ID:              1,
			DefaultCarID:    1,
			Nickname:        "Boss",
			ProfileImageURL: "http://localhost:8080/public/BOSS.PNG",
			CreateTime:      now,
			UpdateTime:      now,
		},
		{
			ID:              2,
			DefaultCarID:    1,
			Nickname:        "Best",
			ProfileImageURL: "http://localhost:8080/public/BEST.PNG",
			CreateTime:      now,
			UpdateTime:      now,
		},
		{
			ID:              3,
			DefaultCarID:    2,
			Nickname:        "Nut",
			ProfileImageURL: "http://localhost:8080/public/NUT.PNG",
			CreateTime:      now,
			UpdateTime:      now,
		},
		{
			ID:              4,
			DefaultCarID:    1,
			Nickname:        "Pat",
			ProfileImageURL: "http://localhost:8080/public/PAT.PNG",
			CreateTime:      now,
			UpdateTime:      now,
		},
	}

	fuelRefills := []domains.FuelRefill{
		{
			ID:                    1,
			CarID:                 1,
			RefillTime:            time.Date(2023, time.January, 1, 0, 0, 0, 0, loc),
			TotalMoney:            decimal.NewFromFloat(500),
			KilometerBeforeRefill: 300,
			KilometerAfterRefill:  800,
			FuelPriceCalculated:   decimal.NewFromFloat(1),
			IsPaid:                true,
			CreateBy:              1,
			CreateTime:            now,
			UpdateBy:              1,
			UpdateTime:            now,
		},
		{
			ID:                    2,
			CarID:                 2,
			RefillTime:            time.Date(2023, time.January, 1, 0, 0, 0, 0, loc),
			TotalMoney:            decimal.NewFromFloat(500),
			KilometerBeforeRefill: 400,
			KilometerAfterRefill:  900,
			FuelPriceCalculated:   decimal.NewFromFloat(1),
			IsPaid:                true,
			CreateBy:              1,
			CreateTime:            now,
			UpdateBy:              1,
			UpdateTime:            now,
		},
	}

	fuelUsages := []domains.FuelUsage{
		{
			ID:                 1,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 800,
			KilometerAfterUse:  700,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(100),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 2,
			CarID:              2,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 900,
			KilometerAfterUse:  800,
			Description:        "mock2",
			TotalMoney:         decimal.NewFromFloat(100),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 3,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 700,
			KilometerAfterUse:  680,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 4,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 680,
			KilometerAfterUse:  660,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 5,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 660,
			KilometerAfterUse:  640,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 6,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 640,
			KilometerAfterUse:  620,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 7,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 620,
			KilometerAfterUse:  600,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 8,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 600,
			KilometerAfterUse:  580,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 9,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 580,
			KilometerAfterUse:  560,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			CreateTime:         now,
			UpdateTime:         now,
		},
		{
			ID:                 10,
			CarID:              1,
			FuelUseTime:        now,
			FuelPrice:          decimal.NewFromFloat(1),
			KilometerBeforeUse: 560,
			KilometerAfterUse:  540,
			Description:        "mock",
			TotalMoney:         decimal.NewFromFloat(20),
			CreateTime:         now,
			UpdateTime:         now,
		},
	}

	fuelUsageUsers := []domains.FuelUsageUser{
		{
			ID: 1, FuelUsageID: 1, UserID: 1, IsPaid: false,
		},
		{
			ID: 2, FuelUsageID: 1, UserID: 2, IsPaid: true,
		},
		{
			ID: 3, FuelUsageID: 2, UserID: 3, IsPaid: false,
		},
		{
			ID: 4, FuelUsageID: 2, UserID: 4, IsPaid: true,
		},
		{
			ID: 5, FuelUsageID: 3, UserID: 1, IsPaid: true,
		},
		{
			ID: 6, FuelUsageID: 4, UserID: 1, IsPaid: true,
		},
		{
			ID: 7, FuelUsageID: 5, UserID: 1, IsPaid: true,
		},
		{
			ID: 8, FuelUsageID: 6, UserID: 1, IsPaid: true,
		},
		{
			ID: 9, FuelUsageID: 7, UserID: 1, IsPaid: true,
		},
		{
			ID: 10, FuelUsageID: 8, UserID: 1, IsPaid: true,
		},
		{
			ID: 11, FuelUsageID: 9, UserID: 1, IsPaid: true,
		},
		{
			ID: 12, FuelUsageID: 10, UserID: 1, IsPaid: true,
		},
	}

	if err := db.Create(&cars).Error; err != nil {
		logger.Error(err.Error())
		return
	}

	if err := db.Create(&users).Error; err != nil {
		logger.Error(err.Error())
		return
	}

	if err := db.Create(&fuelRefills).Error; err != nil {
		logger.Error(err.Error())
		return
	}

	if err := db.Create(&fuelUsages).Error; err != nil {
		logger.Error(err.Error())
		return
	}

	if err := db.Create(&fuelUsageUsers).Error; err != nil {
		logger.Error(err.Error())
		return
	}

	slog.Info("successfully init data")
}
