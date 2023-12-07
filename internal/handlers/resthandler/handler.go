package resthandler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bosskrub9992/fuel-management/internal/models"
	"github.com/bosskrub9992/fuel-management/internal/services"
	"github.com/jinleejun-corp/corelib/errs"
	"github.com/labstack/echo/v4"
)

type RESTHandler struct {
	serverStartTime time.Time
	service         *services.Service
}

func New(service *services.Service) *RESTHandler {
	return &RESTHandler{
		serverStartTime: time.Now(),
		service:         service,
	}
}

func (h RESTHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewUnknown(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	return c.JSON(http.StatusOK, users)
}

func (h RESTHandler) GetFuelUsages(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.GetCarFuelUsagesRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewBind(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	data, err := h.service.GetCarFuelUsages(ctx, req)
	if err != nil {
		response := errs.NewUnknown(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	return c.JSON(http.StatusOK, data)
}

func (h RESTHandler) PostFuelUsages(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateFuelUsageRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewBind(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	if err := h.service.CreateFuelUsage(ctx, req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewUnknown(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h RESTHandler) PutFuelUsage(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.PutFuelUsageRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewBind(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	return c.JSON(200, nil)
}

func (h RESTHandler) DeleteFuelUsages(c echo.Context) error {
	return c.JSON(200, nil)
}

func (h RESTHandler) GetFuelRefills(c echo.Context) error {
	return c.JSON(200, nil)
}

func (h RESTHandler) PostFuelRefills(c echo.Context) error {
	return c.JSON(200, nil)
}

func (h RESTHandler) PutFuelRefill(c echo.Context) error {
	return c.JSON(200, nil)
}

func (h RESTHandler) DeleteFuelRefills(c echo.Context) error {
	return c.JSON(200, nil)
}

func (h RESTHandler) GetHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, h.serverStartTime)
}
