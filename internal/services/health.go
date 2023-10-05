package services

import (
	"context"
	"time"

	"github.com/bosskrub9992/fuel-management/internal/models"
)

type HealthService struct {
	ServerStartTime time.Time
}

func NewHealthService() *HealthService {
	return &HealthService{
		ServerStartTime: time.Now(),
	}
}

func (s *HealthService) GetHealth(ctx context.Context) *models.GetHealthResponse {
	return &models.GetHealthResponse{
		ServerStartTime: s.ServerStartTime,
	}
}
