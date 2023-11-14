package htmxhandler

import "github.com/bosskrub9992/fuel-management/internal/services"

type HTMXHandler struct {
	service *services.Service
}

func New(
	service *services.Service,
) *HTMXHandler {
	return &HTMXHandler{
		service: service,
	}
}
