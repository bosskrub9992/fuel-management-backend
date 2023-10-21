package resthandler

import (
	"net/http"

	"github.com/bosskrub9992/fuel-management/internal/services"
	"github.com/labstack/echo/v4"
)

type RESTHandler struct {
	healthService *services.HealthService
}

func NewRESTHandler(
	healthService *services.HealthService,
) *RESTHandler {
	return &RESTHandler{
		healthService: healthService,
	}
}

func (h RESTHandler) GetHealth(c echo.Context) error {
	ctx := c.Request().Context()
	data := h.healthService.GetHealth(ctx)
	return c.JSON(http.StatusOK, data)
}
