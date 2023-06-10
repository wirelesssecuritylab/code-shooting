package router

import (
	"github.com/labstack/echo/v4"

	"code-shooting/interface/controller"
)

func regECRouters(g *echo.Group) {
	g.POST("/ecs", controller.NewECController().ImportECs)
	g.GET("/ecs", controller.NewECController().QueryECs)
	g.DELETE("/ecs/:id", controller.NewECController().DeleteEC)
}
