package resthandler

import "github.com/bosskrub9992/fuel-management/internal/services"

type RESTHandler struct {
	healthService *services.HealthService
	service       *services.Service
}

func NewRESTHandler(
	healthService *services.HealthService,
	service *services.Service,
) *RESTHandler {
	return &RESTHandler{
		healthService: healthService,
		service:       service,
	}
}
