package routes

import (
	"github.com/gin-gonic/gin"

	"invoice-service/middlewares"

	"invoice-service/controllers"
)

type InvoiceRoute struct {
	controller controllers.IControllerRegistry
	route      *gin.RouterGroup
}

type IInvoiceRoute interface {
	Run()
}

func NewInvoiceRoute(
	controller controllers.IControllerRegistry,
	route *gin.RouterGroup,
) IInvoiceRoute {
	return &InvoiceRoute{
		controller: controller,
		route:      route,
	}
}

func (o *InvoiceRoute) Run() {
	group := o.route.Group("/invoice")
	group.POST("/generate", middlewares.StaticAPIKey(), o.controller.GetInvoice().StoreInvoice)
}
