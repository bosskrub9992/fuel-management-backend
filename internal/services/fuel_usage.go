package services

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/bosskrub9992/fuel-management/internal/domains"
	"github.com/bosskrub9992/fuel-management/internal/models"
	"github.com/shopspring/decimal"
)

func (s *Service) CreateFuelUsage(ctx context.Context, req models.CreateFuelUsageRequest) (*models.GetAllFuelUsageData, error) {
	kmUsed := decimal.NewFromInt(req.KilometerBeforeUse - req.KilometerAfterUse)
	totalMoney := kmUsed.Mul(req.FuelPrice)

	fuelUseDate, err := time.Parse("2006-01-02", req.FuelUseDate)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	newFuelUsage := domains.FuelUsage{
		FuelUseDate:        fuelUseDate,
		FuelPrice:          req.FuelPrice,
		KilometerBeforeUse: req.KilometerBeforeUse,
		KilometerAfterUse:  req.KilometerAfterUse,
		Description:        req.Description,
		TotalMoney:         totalMoney,
		CreateTime:         time.Now(),
		UpdateTime:         time.Now(),
	}

	newFuelUsageID, err := s.db.CreateFuelUsage(ctx, newFuelUsage)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	var newFuelUsageUsers []domains.FuelUsageUser
	for _, userID := range req.UserIDs {
		newFuelUsageUsers = append(newFuelUsageUsers, domains.FuelUsageUser{
			FuelUsageID: newFuelUsageID,
			UserID:      userID,
		})
	}

	if len(newFuelUsageUsers) > 0 {
		if err := s.db.CreateFuelUsageUsers(ctx, newFuelUsageUsers); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, err
		}
	}

	return s.GetAllFuelUsage(ctx)
}

func (s *Service) GetAllFuelUsage(ctx context.Context) (*models.GetAllFuelUsageData, error) {

	latestFuelRefill, err := s.db.GetLatestFuelRefill(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	allFuelUsage, err := s.db.GetAllFuelUsageWithUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	var latestKilometerAfterUse int64
	if len(allFuelUsage) > 0 {
		latestKilometerAfterUse = allFuelUsage[0].KilometerAfterUse
	}

	var getAllFuelUsageData = models.GetAllFuelUsageData{
		TodayDate:               time.Now().Format("2006-01-02"),
		LatestKilometerAfterUse: latestKilometerAfterUse,
		LatestFuelPrice:         latestFuelRefill.FuelPriceCalculated,
		AllUsers:                []models.User{},
		Data:                    []models.GetAllFuelUsageDatum{},
	}

	for _, user := range users {
		getAllFuelUsageData.AllUsers = append(getAllFuelUsageData.AllUsers, models.User{
			ID:       user.ID,
			Nickname: user.Nickname,
		})
	}

	for _, fuelUsage := range allFuelUsage {
		fuelUsageDatum := models.GetAllFuelUsageDatum{
			ID:                 fuelUsage.ID,
			FuelUseDate:        fuelUsage.FuelUseDate.Format("2006-01-02"),
			FuelPrice:          fuelUsage.FuelPrice,
			KilometerBeforeUse: fuelUsage.KilometerBeforeUse,
			KilometerAfterUse:  fuelUsage.KilometerAfterUse,
			Description:        fuelUsage.Description,
			TotalMoney:         fuelUsage.TotalMoney,
			FuelUsers:          strings.Join(fuelUsage.Users, ", "),
		}
		getAllFuelUsageData.Data = append(getAllFuelUsageData.Data, fuelUsageDatum)
	}

	return &getAllFuelUsageData, nil
}
