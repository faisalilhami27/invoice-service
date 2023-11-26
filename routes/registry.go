package routes

import (
	"github.com/gin-gonic/gin"

	controllerRegistry "invoice-service/controllers"
	"invoice-service/middlewares"
	invoiceRoute "invoice-service/routes/invoice"
	templateRoute "invoice-service/routes/template"
)

type IRouteRegistry interface {
	Serve()
}

type Route struct {
	controller controllerRegistry.IControllerRegistry
	Route      *gin.RouterGroup
}

func NewRouteRegistry(
	controller controllerRegistry.IControllerRegistry,
	route *gin.RouterGroup,
) IRouteRegistry {
	return &Route{
		controller: controller,
		Route:      route,
	}
}

func (r *Route) Serve() {
	r.Route.Use(middlewares.HandlePanic)
	r.templateRoute().Run()
	r.invoiceRoute().Run()
}

func (r *Route) templateRoute() templateRoute.ITemplateRoute {
	return templateRoute.NewTemplateRoute(r.controller, r.Route)
}

func (r *Route) invoiceRoute() invoiceRoute.IInvoiceRoute {
	return invoiceRoute.NewInvoiceRoute(r.controller, r.Route)
}
