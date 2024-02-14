package httphandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/bosskrub9992/fuel-management-backend/internal/entities/models"
	"github.com/bosskrub9992/fuel-management-backend/internal/services"
	"github.com/bosskrub9992/fuel-management-backend/library/errs"
)

type HTTPHandler struct {
	service         *services.Service
	serverStartTime time.Time
}

func New(service *services.Service, serverStartTime time.Time) (*HTTPHandler, error) {
	if service == nil {
		return nil, errs.ErrNotEnoughArgForDependencyInjection
	}
	return &HTTPHandler{
		service:         service,
		serverStartTime: serverStartTime,
	}, nil
}

func (h HTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, users)
}

func (h HTTPHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carData, err := h.service.GetCars(ctx)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, carData)
}

func (h HTTPHandler) GetFuelUsages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryParams := r.URL.Query()
	currentCarID, err := strconv.ParseInt(queryParams.Get("currentCarId"), 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}
	currentUserID, err := strconv.ParseInt(queryParams.Get("currentUserId"), 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}
	pageIndex, err := strconv.Atoi(queryParams.Get("pageIndex"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}
	pageSize, err := strconv.Atoi(queryParams.Get("pageSize"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	req := models.GetFuelUsagesRequest{
		CurrentCarID:  currentCarID,
		CurrentUserID: currentUserID,
		PageIndex:     pageIndex,
		PageSize:      pageSize,
	}

	data, err := h.service.GetFuelUsages(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, data)
}

func (h HTTPHandler) PostFuelUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateFuelUsageRequest

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	if err := h.service.CreateFuelUsage(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusCreated, nil)
}

func (h HTTPHandler) PutFuelUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fuelUsageID, err := strconv.ParseInt(r.PathValue("fuelUsageId"), 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	req := models.PutFuelUsageRequest{
		FuelUsageID: fuelUsageID,
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	if err := h.service.UpdateFuelUsage(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, nil)
}

func (h HTTPHandler) DeleteFuelUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fuelUsageID, err := strconv.Atoi(r.PathValue("fuelUsageId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	req := models.DeleteFuelUsageByIDRequest{
		FuelUsageID: int64(fuelUsageID),
	}

	if err := h.service.DeleteFuelUsageByID(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, nil)
}

func (h HTTPHandler) GetFuelUsageByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fuelUsageID, err := strconv.Atoi(r.PathValue("fuelUsageId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	req := models.GetFuelUsageByIDRequest{
		FuelUsageID: int64(fuelUsageID),
	}

	data, err := h.service.GetFuelUsageByID(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, data)
}

func (h HTTPHandler) GetFuelRefills(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.GetFuelRefillRequest

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	data, err := h.service.GetFuelRefills(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, data)
}

func (h HTTPHandler) PostFuelRefill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateFuelRefillRequest

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	if err := h.service.CreateFuelRefill(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusCreated, nil)
}

func (h HTTPHandler) GetFuelRefillByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fuelRefillID, err := strconv.Atoi(r.PathValue("fuelRefillId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	req := models.GetFuelRefillByIDRequest{
		FuelRefillID: int64(fuelRefillID),
	}

	response, err := h.service.GetFuelRefillByID(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, response)
}

func (h HTTPHandler) PutFuelRefillByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fuelRefillID, err := strconv.Atoi(r.PathValue("fuelRefillId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	req := models.PutFuelRefillByIDRequest{
		FuelRefillID: int64(fuelRefillID),
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	if err := h.service.UpdateFuelRefillByID(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, nil)
}

func (h HTTPHandler) DeleteFuelRefillByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fuelRefillID, err := strconv.Atoi(r.PathValue("fuelRefillId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	req := models.DeleteFuelRefillByIDRequest{
		FuelRefillID: int64(fuelRefillID),
	}

	if err := h.service.DeleteFuelRefillByID(ctx, req); err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, nil)
}

func (h HTTPHandler) GetLatestFuelInfoResponse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	carID, err := strconv.Atoi(r.URL.Query().Get("carId"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		sendBadRequest(w, r)
		return
	}

	req := models.GetLatestFuelInfoRequest{
		CarID: int64(carID),
	}

	response, err := h.service.GetLatestFuelInfoResponse(ctx, req)
	if err != nil {
		if response, ok := err.(errs.Err); ok {
			sendJSON(w, r, response.Status, response)
			return
		}
		response := errs.ErrAPIFailed
		sendJSON(w, r, response.Status, response)
		return
	}

	sendJSON(w, r, http.StatusOK, response)
}

func (h HTTPHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, r, http.StatusOK, struct {
		ServerStartTime time.Time `json:"serverStartTime"`
	}{
		ServerStartTime: h.serverStartTime,
	})
}
