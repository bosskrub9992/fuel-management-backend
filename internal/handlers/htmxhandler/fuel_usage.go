package htmxhandler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/bosskrub9992/fuel-management/internal/models"
	"github.com/jinleejun-corp/corelib/errs"
	"github.com/labstack/echo/v4"
)

func (h *HTMXHandler) FuelUsage(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.service.GetAllFuelUsage(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewUnknown(err)
		return c.JSON(response.HTTPStatusCode, response)
	}
	return c.Render(http.StatusOK, "fuel-usage", data)
}

func (h *HTMXHandler) CreateFuelUsage(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateFuelUsageRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewBind(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	params, err := c.FormParams()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewUnknown(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	if userIDs, found := params["userIdCheckbox"]; found {
		for _, rawUserID := range userIDs {
			userID, err := strconv.Atoi(rawUserID)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				response := errs.NewUnknown(err)
				return c.JSON(response.HTTPStatusCode, response)
			}
			req.UserIDs = append(req.UserIDs, int64(userID))
		}
	}

	if err := req.Validate(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewValidate(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	data, err := h.service.CreateFuelUsage(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		response := errs.NewUnknown(err)
		return c.JSON(response.HTTPStatusCode, response)
	}

	c.Response().Header().Set("HX-Trigger-After-Swap", "closeCreateFuelUsageModal")

	return c.Render(http.StatusOK, "create-fuel-usage", data)
}
