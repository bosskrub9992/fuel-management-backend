package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/internal/entities/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/entities/models"
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

func (s *Service) DeleteFuelUsageByID(ctx context.Context, req models.DeleteFuelUsageByIDRequest) error {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return errs.ErrValidateFailed
	}

	return s.db.Transaction(ctx, func(ctxTx context.Context) error {
		if err := s.db.DeleteFuelUsageByID(ctxTx, req.FuelUsageID); err != nil {
			slog.ErrorContext(ctxTx, err.Error())
			return err
		}

		if err := s.db.DeleteFuelUsageUsersByFuelUsageID(ctxTx, req.FuelUsageID); err != nil {
			slog.ErrorContext(ctxTx, err.Error())
			return err
		}

		return nil
	})
}

func (s *Service) UpdateFuelUsage(ctx context.Context, req models.PutFuelUsageRequest) error {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return errs.ErrValidateFailed
	}

	oldfuelUsage, err := s.db.GetFuelUsageByID(ctx, req.FuelUsageID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	totalMoney, err := calculateTotalMoney(
		req.KilometerBeforeUse,
		req.KilometerAfterUse,
		req.FuelPrice,
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return s.db.Transaction(ctx, func(ctxTx context.Context) error {
		fuelUsage := domains.FuelUsage{
			ID:                 req.FuelUsageID,
			CarID:              req.CurrentCarID,
			FuelUseTime:        req.FuelUseTime,
			FuelPrice:          req.FuelPrice,
			KilometerBeforeUse: req.KilometerBeforeUse,
			KilometerAfterUse:  req.KilometerAfterUse,
			Description:        req.Description,
			TotalMoney:         totalMoney,
			CreateTime:         oldfuelUsage.CreateTime,
			UpdateTime:         time.Now(),
		}

		if err := s.db.UpdateFuelUsage(ctxTx, fuelUsage); err != nil {
			slog.ErrorContext(ctxTx, err.Error())
			return err
		}

		if err := s.db.DeleteFuelUsageUsersByFuelUsageID(ctxTx, req.FuelUsageID); err != nil {
			slog.ErrorContext(ctxTx, err.Error())
			return err
		}

		newFuelUsageUsers := []domains.FuelUsageUser{}
		for _, fuelUser := range req.FuelUsers {
			newFuelUsageUsers = append(newFuelUsageUsers, domains.FuelUsageUser{
				FuelUsageID: req.FuelUsageID,
				UserID:      fuelUser.UserID,
				IsPaid:      fuelUser.IsPaid,
			})
		}

		if err := s.db.CreateFuelUsageUsers(ctxTx, newFuelUsageUsers); err != nil {
			slog.ErrorContext(ctxTx, err.Error())
			return err
		}

		return nil
	})
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
		return nil, errs.ErrValidateFailed
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
		CarID:     req.CurrentCarID,
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
		Now:                     time.Now(),
		LatestKilometerAfterUse: latestKilometerAfterUse,
		LatestFuelPrice:         latestFuelRefill.FuelPriceCalculated,
		AllUsers:                allUsers,
		CurrentUser:             currentUser,
		Data: func() (data []models.CarFuelUsageDatum) {
			for _, fuelUsage := range carFuelUsages {
				fuelUsageDatum := models.CarFuelUsageDatum{
					ID:                 fuelUsage.ID,
					FuelUseTime:        fuelUsage.FuelUseTime.Format("2006-01-02"),
					FuelPrice:          fuelUsage.FuelPrice,
					KilometerBeforeUse: fuelUsage.KilometerBeforeUse,
					KilometerAfterUse:  fuelUsage.KilometerAfterUse,
					Description:        fuelUsage.Description,
					TotalMoney:         fuelUsage.TotalMoney,
					FuelUsers:          strings.Join(fuelUsage.Users, ", ")}
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
				if car.ID == req.CurrentCarID {
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
		return errs.ErrValidateFailed
	}

	totalMoney, err := calculateTotalMoney(
		req.KilometerBeforeUse,
		req.KilometerAfterUse,
		req.FuelPrice,
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	fuelUsage := domains.FuelUsage{
		CarID:              req.CurrentCarID,
		FuelUseTime:        req.FuelUseTime,
		FuelPrice:          req.FuelPrice,
		KilometerBeforeUse: req.KilometerBeforeUse,
		KilometerAfterUse:  req.KilometerAfterUse,
		Description:        req.Description,
		TotalMoney:         totalMoney,
		CreateTime:         time.Now(),
		UpdateTime:         time.Now(),
	}

	return s.db.Transaction(ctx, func(ctxTx context.Context) error {
		fuelUsageID, err := s.db.CreateFuelUsage(ctxTx, fuelUsage)
		if err != nil {
			slog.ErrorContext(ctxTx, err.Error())
			return err
		}

		var fuelUsageUsers []domains.FuelUsageUser
		for _, fuelUser := range req.FuelUsers {
			fuelUsageUsers = append(fuelUsageUsers, domains.FuelUsageUser{
				FuelUsageID: fuelUsageID,
				UserID:      fuelUser.UserID,
				IsPaid:      fuelUser.IsPaid,
			})
		}

		if err := s.db.CreateFuelUsageUsers(ctxTx, fuelUsageUsers); err != nil {
			slog.ErrorContext(ctxTx, err.Error())
			return err
		}

		return nil
	})
}

func (s *Service) GetFuelUsageByID(ctx context.Context, req models.GetFuelUsageByIDRequest) (*models.GetFuelUsageByIDResponse, error) {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errs.ErrValidateFailed
	}

	fuelUsage, err := s.db.GetFuelUsageByID(ctx, req.FuelUsageID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	fuelUsageUsers, err := s.db.GetFuelUsageUsersByFuelUsageID(ctx, fuelUsage.ID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	fuelUsers := []models.GetFuelUser{}
	for _, fuelUsageUser := range fuelUsageUsers {
		fuelUsers = append(fuelUsers, models.GetFuelUser{
			UserID:   fuelUsageUser.UserID,
			Nickname: fuelUsageUser.Nickname,
			IsPaid:   fuelUsageUser.IsPaid,
		})
	}

	response := models.GetFuelUsageByIDResponse{
		FuelUseTime:        fuelUsage.FuelUseTime,
		FuelPrice:          fuelUsage.FuelPrice,
		FuelUsers:          fuelUsers,
		Description:        fuelUsage.Description,
		KilometerBeforeUse: fuelUsage.KilometerBeforeUse,
		KilometerAfterUse:  fuelUsage.KilometerAfterUse,
		TotalMoney:         fuelUsage.TotalMoney,
		EachShouldPay:      fuelUsage.TotalMoney.DivRound(decimal.NewFromInt(int64(len(fuelUsageUsers))), 2),
	}

	return &response, nil
}

func calculateTotalMoney(kmBeforeUse, kmAfterUse int64, fuelPrice decimal.Decimal) (decimal.Decimal, error) {
	if kmBeforeUse < kmAfterUse {
		return decimal.Zero, fmt.Errorf(
			"kmBeforeUse should be > kmAfterUse, kmBeforeUse: [%d], kmAfterUse: [%d]",
			kmBeforeUse,
			kmAfterUse,
		)
	}
	kmUsed := decimal.NewFromInt(kmBeforeUse - kmAfterUse)
	return kmUsed.Mul(fuelPrice), nil
}
