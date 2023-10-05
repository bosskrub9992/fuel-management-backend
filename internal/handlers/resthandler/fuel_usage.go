package resthandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *RESTHandler) FuelUsage(c echo.Context) error {
	return c.Render(http.StatusOK, "fuel-usage", "fuel-usage")
}
