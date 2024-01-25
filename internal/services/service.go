package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/config"
	"github.com/bosskrub9992/fuel-management-backend/internal/domains"
	"github.com/bosskrub9992/fuel-management-backend/internal/models"
	"github.com/bosskrub9992/fuel-management-backend/library/errs"
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
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	var userData []models.GetUserDatum
	for _, user := range users {
		userData = append(userData, models.GetUserDatum{
			ID:              user.ID,
			DefaultCarID:    user.DefaultCarID,
			Nickname:        user.Nickname,
			ProfileImageURL: user.ProfileImageURL,
		})
	}

	return &models.GetUserData{
		Data: userData,
	}, nil
}

func (s *Service) GetCars(ctx context.Context) (*models.GetCarData, error) {
	cars, err := s.db.GetAllCars(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	var carData []models.CarDatum
	for _, car := range cars {
		carData = append(carData, models.CarDatum{
			ID:   car.ID,
			Name: car.Name,
		})
	}

	return &models.GetCarData{
		Data: carData,
	}, nil
}

func (s *Service) GetFuelUsages(ctx context.Context, req models.GetFuelUsagesRequest) (*models.GetFuelUsagesResponse, error) {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errs.ErrValidateFailed
	}

	fuelUsages, totalRecord, err := s.db.GetFuelUsageInPagination(ctx, GetFuelUsageInPaginationParams{
		CarID:     req.CurrentCarID,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	fuelUsageIDs := []int64{}
	for _, fuelUsage := range fuelUsages {
		fuelUsageIDs = append(fuelUsageIDs, fuelUsage.ID)
	}

	fuelUsers, err := s.db.GetFuelUsageUsersByFuelUsageIDs(ctx, fuelUsageIDs)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	usageIdToUser := make(map[int64][]FuelUsageUser)
	for _, user := range fuelUsers {
		_, found := usageIdToUser[user.FuelUsageID]
		if !found {
			usageIdToUser[user.FuelUsageID] = []FuelUsageUser{}
		}
		usageIdToUser[user.FuelUsageID] = append(usageIdToUser[user.FuelUsageID], user)
	}

	var fuelUsageData []models.FuelUsageDatum
	for _, fuelUsage := range fuelUsages {
		var fuelUserNames []string
		fuelUsers, found := usageIdToUser[fuelUsage.ID]
		if !found {
			err := fmt.Errorf("not found fuel usage id: [%d]", fuelUsage.ID)
			slog.ErrorContext(ctx, err.Error())
			return nil, err
		}
		for _, fuelUser := range fuelUsers {
			isPaid := "❌"
			if fuelUser.IsPaid {
				isPaid = "✅"
			}
			fuelUserNames = append(fuelUserNames, fmt.Sprintf("%s%s",
				isPaid,
				fuelUser.Nickname,
			))
		}
		fuelUsageData = append(fuelUsageData, models.FuelUsageDatum{
			ID:                 fuelUsage.ID,
			FuelUseTime:        fuelUsage.FuelUseTime.Format("_2 Jan 15:04"),
			FuelPrice:          fuelUsage.FuelPrice,
			KilometerBeforeUse: fuelUsage.KilometerBeforeUse,
			KilometerAfterUse:  fuelUsage.KilometerAfterUse,
			Description:        fuelUsage.Description,
			TotalMoney:         fuelUsage.TotalMoney,
			FuelUsers:          strings.Join(fuelUserNames, " "),
		})
	}

	return &models.GetFuelUsagesResponse{
		FuelUsageData: fuelUsageData,
		TotalRecord:   totalRecord,
		TotalPage:     int64(math.Ceil(float64(totalRecord) / float64(req.PageSize))),
	}, nil
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

func (s *Service) GetFuelRefills(ctx context.Context, req models.GetFuelRefillRequest) (*models.GetFuelRefillResponse, error) {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errs.ErrValidateFailed
	}

	fuelRefills, totalRecord, err := s.db.GetFuelRefillPagination(ctx, GetFuelRefillPaginationParams{
		CarID:     req.CurrentCarID,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	response := models.GetFuelRefillResponse{
		FuelRefillData: []models.FuelRefillDatum{},
		TotalRecord:    totalRecord,
		TotalPage:      int(math.Ceil(float64(totalRecord) / float64(req.PageSize))),
	}

	for _, fr := range fuelRefills {
		response.FuelRefillData = append(response.FuelRefillData, models.FuelRefillDatum{
			ID:                    fr.ID,
			RefillTime:            fr.RefillTime.Format("_2 Jan 15:04"),
			KilometerBeforeRefill: fr.KilometerBeforeRefill,
			KilometerAfterRefill:  fr.KilometerAfterRefill,
			TotalMoney:            fr.TotalMoney,
			FuelPriceCalculated:   fr.FuelPriceCalculated,
			IsPaid:                fr.IsPaid,
		})
	}

	return &response, nil
}

func (s *Service) CreateFuelRefill(ctx context.Context, req models.CreateFuelRefillRequest) error {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return errs.ErrValidateFailed
	}

	fuelPrice, err := calculateFuelPrice(
		req.TotalMoney,
		req.KilometerBeforeRefill,
		req.KilometerAfterRefill,
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	now := time.Now()

	fuelRefill := domains.FuelRefill{
		CarID:                 req.CurrentCarID,
		RefillTime:            req.RefillTime,
		TotalMoney:            req.TotalMoney,
		KilometerBeforeRefill: req.KilometerBeforeRefill,
		KilometerAfterRefill:  req.KilometerAfterRefill,
		FuelPriceCalculated:   fuelPrice,
		IsPaid:                req.IsPaid,
		UpdateBy:              req.CurrentUserID,
		CreateBy:              req.CurrentUserID,
		CreateTime:            now,
		UpdateTime:            now,
	}

	if err := s.db.CreateFuelRefill(ctx, fuelRefill); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}

func (s *Service) GetFuelRefillByID(ctx context.Context, req models.GetFuelRefillByIDRequest) (*models.GetFuelRefillByIDResponse, error) {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errs.ErrValidateFailed
	}

	fuelRefill, err := s.db.GetFuelRefillByID(ctx, req.FuelRefillID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	return &models.GetFuelRefillByIDResponse{
		RefillTime:            fuelRefill.RefillTime,
		KilometerBeforeRefill: fuelRefill.KilometerBeforeRefill,
		KilometerAfterRefill:  fuelRefill.KilometerAfterRefill,
		TotalMoney:            fuelRefill.TotalMoney,
		FuelPriceCalculated:   fuelRefill.FuelPriceCalculated,
		IsPaid:                fuelRefill.IsPaid,
	}, nil
}

func (s *Service) UpdateFuelRefillByID(ctx context.Context, req models.PutFuelRefillByIDRequest) error {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return errs.ErrValidateFailed
	}

	oldFuelRefill, err := s.db.GetFuelRefillByID(ctx, req.FuelRefillID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	newFuelPrice, err := calculateFuelPrice(
		req.TotalMoney,
		req.KilometerBeforeRefill,
		req.KilometerAfterRefill,
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	newFuelRefill := domains.FuelRefill{
		ID:                    req.FuelRefillID,
		CarID:                 req.CurrentCarID,
		RefillTime:            req.RefillTime,
		TotalMoney:            req.TotalMoney,
		KilometerBeforeRefill: req.KilometerBeforeRefill,
		KilometerAfterRefill:  req.KilometerAfterRefill,
		FuelPriceCalculated:   newFuelPrice,
		IsPaid:                req.IsPaid,
		CreateBy:              oldFuelRefill.CreateBy,
		CreateTime:            oldFuelRefill.CreateTime,
		UpdateBy:              req.CurrentUserID,
		UpdateTime:            time.Now(),
	}

	if err := s.db.UpdateFuelRefill(ctx, newFuelRefill); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}

func calculateFuelPrice(
	totalMoney decimal.Decimal,
	kmBeforeRefill, kmAfterRefill int64,
) (
	decimal.Decimal, error,
) {
	if kmAfterRefill <= kmBeforeRefill {
		return decimal.Zero, errors.New("kmAfterRefill should > kmBeforeRefill")
	}
	increaseKm := decimal.NewFromInt(kmAfterRefill - kmBeforeRefill)
	return totalMoney.DivRound(increaseKm, 2), nil
}

func (s *Service) DeleteFuelRefillByID(ctx context.Context, req models.DeleteFuelRefillByIDRequest) error {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return errs.ErrValidateFailed
	}

	if err := s.db.DeleteFuelRefillByID(ctx, req.FuelRefillID); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}

func (s *Service) GetLatestFuelInfoResponse(ctx context.Context, req models.GetLatestFuelInfoRequest) (*models.GetLatestFuelInfoResponse, error) {
	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, errs.ErrValidateFailed
	}

	latestFuelUsage, err := s.db.GetLatestFuelUsageByCarID(ctx, req.CarID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	latestFuelRefill, err := s.db.GetLatestFuelRefillByCarID(ctx, req.CarID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	var latestKmAfterUse = latestFuelUsage.KilometerAfterUse
	if latestFuelRefill.RefillTime.After(latestFuelUsage.FuelUseTime) {
		latestKmAfterUse = latestFuelRefill.KilometerAfterRefill
	}

	return &models.GetLatestFuelInfoResponse{
		LatestFuelPrice:         latestFuelRefill.FuelPriceCalculated,
		LatestKilometerAfterUse: latestKmAfterUse,
	}, nil
}
