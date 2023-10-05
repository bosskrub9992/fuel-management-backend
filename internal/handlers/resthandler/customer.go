package resthandler

import (
	"net/http"

	"github.com/bosskrub9992/fuel-management/internal/models"

	"github.com/jinleejun-corp/corelib/errs"
	"github.com/jinleejun-corp/corelib/slogger"
	"github.com/labstack/echo/v4"
)

func (h *RESTHandler) CreateCustomer(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
		slogger.Error(ctx, err.Error())
		response := errs.NewBind(err)
		return c.JSON(response.Status, response)
	}

	if err := req.Vaildate(); err != nil {
		slogger.Error(ctx, err.Error())
		response := errs.NewValidate(err)
		return c.JSON(response.Status, response)
	}

	response, err := h.service.CreateCustomer(ctx, req)
	if err != nil {
		response := errs.NewUnknown(err)
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, response)
}
