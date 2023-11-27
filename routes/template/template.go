package routes

import (
	"github.com/gin-gonic/gin"

	"invoice-service/controllers"
)

type TemplateRoute struct {
	controller controllers.IControllerRegistry
	route      *gin.RouterGroup
}

type ITemplateRoute interface {
	Run()
}

func NewTemplateRoute(
	controller controllers.IControllerRegistry,
	route *gin.RouterGroup,
) ITemplateRoute {
	return &TemplateRoute{
		controller: controller,
		route:      route,
	}
}

func (o *TemplateRoute) Run() {
	group := o.route.Group("/template")
	group.GET("", o.controller.GetTemplate().GetTemplate)
	group.POST("/upload", o.controller.GetTemplate().StoreTemplate)
}
