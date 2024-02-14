package routers

import (
	"net/http"

	"github.com/bosskrub9992/fuel-management-backend/internal/handlers/httphandler"
	"github.com/bosskrub9992/fuel-management-backend/library/errs"
)

type HTTPRouter struct {
	mux         *http.ServeMux
	httpHandler *httphandler.HTTPHandler
}

func NewHTTPRouter(mux *http.ServeMux, httpHandler *httphandler.HTTPHandler) (*HTTPRouter, error) {
	if mux == nil || httpHandler == nil {
		return nil, errs.ErrNotEnoughArgForDependencyInjection
	}
	return &HTTPRouter{
		mux:         mux,
		httpHandler: httpHandler,
	}, nil
}

func (r HTTPRouter) Init() *http.ServeMux {
	fs := http.FileServer(http.Dir("public"))
	r.mux.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.mux.HandleFunc("GET /api/v1/cars", r.httpHandler.GetCars)
	r.mux.HandleFunc("GET /api/v1/users", r.httpHandler.GetUsers)

	r.mux.HandleFunc("POST /api/v1/fuel/usages", r.httpHandler.PostFuelUsage)
	r.mux.HandleFunc("GET /api/v1/fuel/usages", r.httpHandler.GetFuelUsages)
	r.mux.HandleFunc("GET /api/v1/fuel/usages/{fuelUsageId}", r.httpHandler.GetFuelUsageByID)
	r.mux.HandleFunc("PUT /api/v1/fuel/usages/{fuelUsageId}", r.httpHandler.PutFuelUsage)
	r.mux.HandleFunc("DELETE /api/v1/fuel/usages/{fuelUsageId}", r.httpHandler.DeleteFuelUsage)

	r.mux.HandleFunc("POST /api/v1/fuel/refills", r.httpHandler.PostFuelRefill)
	r.mux.HandleFunc("GET /api/v1/fuel/refills", r.httpHandler.GetFuelRefills)
	r.mux.HandleFunc("GET /api/v1/fuel/refills/{fuelRefillId}", r.httpHandler.GetFuelRefillByID)
	r.mux.HandleFunc("PUT /api/v1/fuel/refills/{fuelRefillId}", r.httpHandler.PutFuelRefillByID)
	r.mux.HandleFunc("DELETE /api/v1/fuel/refills/{fuelRefillId}", r.httpHandler.DeleteFuelRefillByID)

	r.mux.HandleFunc("GET /latest-fuel-info", r.httpHandler.GetLatestFuelInfoResponse)
	return r.mux
}
