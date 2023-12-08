package resthandler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/models"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
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
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, users)
}

func (h RESTHandler) GetFuelUsages(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.GetCarFuelUsagesRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	data, err := h.service.GetCarFuelUsages(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, data)
}

func (h RESTHandler) PostFuelUsage(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateFuelUsageRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	if err := h.service.CreateFuelUsage(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h RESTHandler) PutFuelUsage(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.PutFuelUsageRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	if err := h.service.UpdateFuelUsage(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h RESTHandler) DeleteFuelUsage(c echo.Context) error {
	ctx := c.Request().Context()

	fuelUsageID, err := strconv.Atoi(c.Param("fuelUsageId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	req := models.DeleteFuelUsageByIDRequest{
		FuelUsageID: int64(fuelUsageID),
	}

	if err := h.service.DeleteFuelUsageByID(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h RESTHandler) GetFuelUsageByID(c echo.Context) error {
	ctx := c.Request().Context()

	fuelUsageID, err := strconv.Atoi(c.Param("fuelUsageId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	req := models.GetFuelUsageByIDRequest{
		FuelUsageID: int64(fuelUsageID),
	}

	data, err := h.service.GetFuelUsageByID(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, data)
}

func (h RESTHandler) GetFuelRefills(c echo.Context) error {
	ctx := c.Request().Context()

	req := models.GetFuelRefillRequest{}
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	data, err := h.service.GetFuelRefills(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, data)
}

func (h RESTHandler) CreateFuelRefill(c echo.Context) error {
	ctx := c.Request().Context()

	req := models.CreateFuelRefillRequest{}
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	if err := h.service.CreateFuelRefill(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h RESTHandler) GetFuelRefillByID(c echo.Context) error {
	ctx := c.Request().Context()

	fuelRefillID, err := strconv.Atoi(c.Param("fuelRefillId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	req := models.GetFuelRefillByIDRequest{
		FuelRefillID: int64(fuelRefillID),
	}

	response, err := h.service.GetFuelRefillByID(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, response)
}

func (h RESTHandler) PutFuelRefillByID(c echo.Context) error {
	ctx := c.Request().Context()

	req := models.PutFuelRefillByIDRequest{}
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	if err := h.service.UpdateFuelRefillByID(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h RESTHandler) DeleteFuelRefillByID(c echo.Context) error {
	ctx := c.Request().Context()

	fuelRefillID, err := strconv.Atoi(c.Param("fuelRefillId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	req := models.DeleteFuelRefillByIDRequest{
		FuelRefillID: int64(fuelRefillID),
	}

	if err := h.service.DeleteFuelRefillByID(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h RESTHandler) GetHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		ServerStartTime time.Time `json:"serverStartTime"`
	}{
		ServerStartTime: h.serverStartTime,
	})
}
