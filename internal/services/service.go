package services

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/internal/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/models"
	"github.com/jinleejun-corp/corelib/errs"
	"github.com/shopspring/decimal"
)

type Service struct {
	cfg *config.Config
	db  DatabaseAdaptor
}

func New(cfg *config.Config, db DatabaseAdaptor) *Service {
	return &Service{
		cfg: cfg,
		db:  db,
	}
}

func (s *Service) GetUsers(ctx context.Context) (*models.GetUserData, error) {
	allUsers, err := s.db.GetAllUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	data := models.GetUserData{
		Data: []models.GetUserDatum{},
	}

	for _, user := range allUsers {
		data.Data = append(data.Data, models.GetUserDatum{
			ID:              user.ID,
			DefaultCarID:    user.DefaultCarID,
			Nickname:        user.Nickname,
			ProfileImageURL: user.ProfileImageURL,
		})
	}

	return &data, nil
}

func (s *Service) GetCarFuelUsages(ctx context.Context, req models.GetCarFuelUsagesRequest) (*models.GetCarFuelUsageData, error) {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewValidate(err)
		return nil, response
	}

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

	getAllFuelUsageData := models.GetCarFuelUsageData{
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

func (s *Service) CreateFuelUsage(ctx context.Context, req models.CreateFuelUsageRequest) error {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewValidate(err)
		return response
	}

	kmUsed := decimal.NewFromInt(req.KilometerBeforeUse - req.KilometerAfterUse)
	totalMoney := kmUsed.Mul(req.FuelPrice)

	newFuelUsage := domains.FuelUsage{
		FuelUseDate:        req.FuelUseDate,
		FuelPrice:          req.FuelPrice,
		KilometerBeforeUse: req.KilometerBeforeUse,
		KilometerAfterUse:  req.KilometerAfterUse,
		Description:        req.Description,
		TotalMoney:         totalMoney,
		CreateTime:         time.Now(),
		UpdateTime:         time.Now(),
	}

	if err := s.db.CreateFuelUsage(ctx, newFuelUsage, req.UserIDs); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}
