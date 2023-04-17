package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/vnnyx/rekadigital-tech-test/internal/delivery/http"
)

type Route struct {
	transactionHandler *http.TransactionHandler
	route              *echo.Echo
}

func NewRoute(transactionHandler *http.TransactionHandler, route *echo.Echo) *Route {
	return &Route{
		transactionHandler: transactionHandler,
		route:              route,
	}
}

func (r *Route) InitRoute() {
	api := r.route.Group("/rekadigital-api")
	api.POST("/transaction", r.transactionHandler.CreateTransaction)
	api.GET("/transaction", r.transactionHandler.GetAllTransaction)
}
