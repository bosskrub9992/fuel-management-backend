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

func (s *Service) GetUsers(ctx context.Context) ([]models.GetUserDatum, error) {
	allUsers, err := s.db.GetAllUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}
	var users []models.GetUserDatum
	for _, user := range allUsers {
		users = append(users, models.GetUserDatum{
			ID:              user.ID,
			DefaultCarID:    user.DefaultCarID,
			Nickname:        user.Nickname,
			ProfileImageURL: user.ProfileImageURL,
		})
	}
	return users, nil
}

func (s *Service) GetCarFuelUsages(ctx context.Context, req models.GetCarFuelUsagesRequest) (*models.GetCarFuelUsageData, error) {
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

	allCars, err := s.db.GetAllCars(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	carFuelUsages, totalRecord, err := s.db.GetCarFuelUsageWithUsers(ctx, GetCarFuelUsageWithUsersParams{
		CarID:     req.CarID,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	var latestKilometerAfterUse int64
	if len(carFuelUsages) > 0 {
		latestKilometerAfterUse = carFuelUsages[0].KilometerAfterUse
	}

	var (
		allUsers    []models.User
		currentUser models.UserWithImageURL
	)
	for _, user := range users {
		allUsers = append(allUsers, models.User{
			ID:       user.ID,
			Nickname: user.Nickname,
		})
		if user.ID == req.CurrentUserID {
			currentUser = models.UserWithImageURL{
				User: models.User{
					ID:       user.ID,
					Nickname: user.Nickname,
				},
				UserImageURL: user.ProfileImageURL,
			}
		}
	}

	var getAllFuelUsageData = models.GetCarFuelUsageData{
		TodayDate:               time.Now().Format("2006-01-02"),
		LatestKilometerAfterUse: latestKilometerAfterUse,
		LatestFuelPrice:         latestFuelRefill.FuelPriceCalculated,
		AllUsers:                allUsers,
		CurrentUser:             currentUser,
		Data: func() (data []models.CarFuelUsageDatum) {
			for _, fuelUsage := range carFuelUsages {
				fuelUsageDatum := models.CarFuelUsageDatum{ID: fuelUsage.ID, FuelUseDate: fuelUsage.FuelUseDate.Format("2006-01-02"), FuelPrice: fuelUsage.FuelPrice, KilometerBeforeUse: fuelUsage.KilometerBeforeUse, KilometerAfterUse: fuelUsage.KilometerAfterUse, Description: fuelUsage.Description, TotalMoney: fuelUsage.TotalMoney, FuelUsers: strings.Join(fuelUsage.Users, ", ")}
				data = append(data, fuelUsageDatum)
			}
			return data
		}(),
		AllCars: func() (cars []models.Car) {
			for _, car := range allCars {
				cars = append(cars, models.Car{ID: car.ID, Name: car.Name})
			}
			return cars
		}(),
		TotalRecord: totalRecord,
		CurrentCar: func() models.Car {
			for _, car := range allCars {
				if car.ID == req.CarID {
					return models.Car{ID: car.ID, Name: car.Name}
				}
			}
			return models.Car{ID: allCars[0].ID, Name: allCars[0].Name}
		}(),
	}

	return &getAllFuelUsageData, nil
}

func (s *Service) CreateFuelUsage(ctx context.Context, req models.CreateFuelUsageRequest) (*models.GetCarFuelUsageData, error) {
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

	return s.GetCarFuelUsages(ctx, models.GetCarFuelUsagesRequest{
		CarID:     req.CarID,
		PageIndex: 1,
		PageSize:  20,
	})
}
