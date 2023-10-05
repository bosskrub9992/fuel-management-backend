package resthandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h RESTHandler) GetHealth(c echo.Context) error {
	ctx := c.Request().Context()
	data := h.healthService.GetHealth(ctx)
	return c.JSON(http.StatusOK, data)
}
