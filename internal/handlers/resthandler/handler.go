package resthandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/models"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
	"github.com/bosskrub9992/fuel-management-backend/library/errs"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type RESTHandler struct {
	service         *services.Service
	serverStartTime time.Time
}

func New(service *services.Service, serverStartTime time.Time) *RESTHandler {
	return &RESTHandler{
		service:         service,
		serverStartTime: serverStartTime,
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

func (h RESTHandler) GetCars(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.service.GetCars(ctx)
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

	var req models.GetFuelUsagesRequest
	if err := c.Bind(&req); err != nil {
		log.Ctx(ctx).Err(err).Send()
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	data, err := h.service.GetFuelUsages(ctx, req)
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
		log.Ctx(ctx).Err(err).Send()
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
		log.Ctx(ctx).Err(err).Send()
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
		log.Ctx(ctx).Err(err).Send()
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
		log.Ctx(ctx).Err(err).Send()
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
		log.Ctx(ctx).Err(err).Send()
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

func (h RESTHandler) PostFuelRefill(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateFuelRefillRequest
	if err := c.Bind(&req); err != nil {
		log.Ctx(ctx).Err(err).Send()
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
		log.Ctx(ctx).Err(err).Send()
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
		log.Ctx(ctx).Err(err).Send()
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
		log.Ctx(ctx).Err(err).Send()
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

func (h RESTHandler) GetLatestFuelInfoResponse(c echo.Context) error {
	ctx := c.Request().Context()

	carID, err := strconv.Atoi(c.QueryParam("carId"))
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	req := models.GetLatestFuelInfoRequest{
		CarID: int64(carID),
	}

	response, err := h.service.GetLatestFuelInfoResponse(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, response)
}

func (h RESTHandler) GetHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		ServerStartTime time.Time `json:"serverStartTime"`
	}{
		ServerStartTime: h.serverStartTime,
	})
}

func (h RESTHandler) GetUserFuelUsages(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	isPaid, err := strconv.ParseBool(c.QueryParam("isPaid"))
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	req := models.GetUserFuelUsagesRequest{
		UserID: int64(userID),
		IsPaid: isPaid,
	}

	data, err := h.service.GetUserFuelUsages(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, data)
}

func (h RESTHandler) GetUserCarUnpaidActivities(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	carID, err := strconv.Atoi(c.Param("carId"))
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	req := models.GetUserCarUnpaidActivitiesRequest{
		UserID: int64(userID),
		CarID:  int64(carID),
	}

	data, err := h.service.GetUserCarUnpaidActivities(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, data)
}

func (h RESTHandler) BulkUpdateUserFuelUsagePaymentStatus(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.BulkUpdateUserFuelUsagePaymentStatusRequest
	if err := c.Bind(&req); err != nil {
		log.Ctx(ctx).Err(err).Send()
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		response := errs.ErrBadRequest
		return c.JSON(response.Status, response)
	}

	req.UserID = int64(userID)

	if err := h.service.BulkUpdateUserFuelUsagePaymentStatus(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h RESTHandler) PayUserCarUnpaidActivities(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.PayUserCarUnpaidActivitiesRequest
	if err := c.Bind(&req); err != nil {
		log.Ctx(ctx).Err(err).Send()
		res := errs.ErrBadRequest
		return c.JSON(res.Status, res)
	}

	if err := h.service.PayUserCarUnpaidActivities(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			return c.JSON(response.Status, response)
		}
		response := errs.ErrAPIFailed
		return c.JSON(response.Status, response)
	}

	return c.JSON(http.StatusOK, nil)
}
