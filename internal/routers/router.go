package routers

import (
	"github.com/bosskrub9992/fuel-management/internal/handlers/resthandler"
	"github.com/jinleejun-corp/corelib/slogger"
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
	r.e.GET("/health", r.restHandler.GetHealth)
	r.e.Use(
		middleware.Recover(),
		middleware.CORS(),
		slogger.MiddlewareREST(),
	)
	r.e.Static("/public", "./public")
	apiV1 := r.e.Group("/api/v1")
	apiV1.GET("/users", r.restHandler.GetUsers)
	apiV1.GET("/fuel/usages", r.restHandler.GetFuelUsages)
	apiV1.POST("/fuel/usages", r.restHandler.PostFuelUsages)
	return r.e
}
