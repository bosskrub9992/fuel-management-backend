package routers

import (
	"net/http"

	"github.com/bosskrub9992/fuel-management-backend/internal/handlers/httphandler"
	"github.com/bosskrub9992/fuel-management-backend/library/errs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type CHIRouter struct {
	mux     *chi.Mux
	handler *httphandler.HTTPHandler
}

func NewCHIRouter(mux *chi.Mux, handler *httphandler.HTTPHandler) (*CHIRouter, error) {
	if mux == nil || handler == nil {
		return nil, errs.ErrNotEnoughArgForDependencyInjection
	}
	return &CHIRouter{
		mux:     mux,
		handler: handler,
	}, nil
}

func (cr CHIRouter) Init() *chi.Mux {
	fs := http.FileServer(http.Dir("public"))
	cr.mux.Handle("/static/*", http.StripPrefix("/static/", fs))

	cr.mux.Use(
		middleware.Recoverer,
		cors.AllowAll().Handler,
	)

	cr.mux.Route("/api/v1", func(r chi.Router) {
		cr.mux.Get("/cars", cr.handler.GetCars)
		cr.mux.Get("/users", cr.handler.GetUsers)

		cr.mux.Post("/fuel/usages", cr.handler.PostFuelUsage)
		cr.mux.Get("/fuel/usages", cr.handler.GetFuelUsages)
		cr.mux.Get("/fuel/usages/{fuelUsageId}", cr.handler.GetFuelUsageByID)
		cr.mux.Put("/fuel/usages/{fuelUsageId}", cr.handler.PutFuelUsage)
		cr.mux.Delete("/fuel/usages/{fuelUsageId}", cr.handler.DeleteFuelUsage)

		cr.mux.Post("/fuel/refills", cr.handler.PostFuelRefill)
		cr.mux.Get("/fuel/refills", cr.handler.GetFuelRefills)
		cr.mux.Get("/fuel/refills/{fuelRefillId}", cr.handler.GetFuelRefillByID)
		cr.mux.Put("/fuel/refills/{fuelRefillId}", cr.handler.PutFuelRefillByID)
		cr.mux.Delete("/fuel/refills/{fuelRefillId}", cr.handler.DeleteFuelRefillByID)

		cr.mux.Get("/latest-fuel-info", cr.handler.GetLatestFuelInfoResponse)
	})

	return cr.mux
}
