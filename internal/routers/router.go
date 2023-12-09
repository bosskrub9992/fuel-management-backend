package routers

import (
	"github.com/bosskrub9992/fuel-management-backend/internal/handlers/resthandler"
	"github.com/jinleejun-corp/corelib/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	e           *echo.Echo
	restHandler *resthandler.RESTHandler
}

func New(e *echo.Echo, restHandler *resthandler.RESTHandler) *Router {
	return &Router{
		e:           e,
		restHandler: restHandler,
	}
}

func (r Router) Init() *echo.Echo {
	r.e.Static("/public", "./public")
	r.e.Use(
		middleware.Recover(),
		middleware.CORS(),
		middlewares.RequestID(),
		middlewares.Logger(),
	)
	r.e.GET("/health", r.restHandler.GetHealth)
	apiV1 := r.e.Group("/api/v1")
	apiV1.GET("/users", r.restHandler.GetUsers)

	apiV1.POST("/fuel/usages", r.restHandler.PostFuelUsage)
	apiV1.GET("/fuel/usages", r.restHandler.GetFuelUsages)
	apiV1.GET("/fuel/usages/:fuelUsageId", r.restHandler.GetFuelUsageByID)
	apiV1.PUT("/fuel/usages/:fuelUsageId", r.restHandler.PutFuelUsage)
	apiV1.DELETE("/fuel/usages/:fuelUsageId", r.restHandler.DeleteFuelUsage)

	apiV1.POST("/fuel/refills", r.restHandler.CreateFuelRefill)
	apiV1.GET("/fuel/refills", r.restHandler.GetFuelRefills)
	apiV1.GET("/fuel/refills/:fuelRefillId", r.restHandler.GetFuelRefillByID)
	apiV1.PUT("/fuel/refills/:fuelRefillId", r.restHandler.PutFuelRefillByID)
	apiV1.DELETE("/fuel/refills/:fuelRefillId", r.restHandler.DeleteFuelRefillByID)
	return r.e
}
