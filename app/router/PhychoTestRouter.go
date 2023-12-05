package router

import (
	"main/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRouter_PhychoTest(controller *controller.PhychoTestController, route *gin.Engine) {

	route.GET("/type_create", controller.CreateType)
	route.POST("/type_addcount", controller.Type_AddCount)
	route.GET("/type_percentage", controller.CountTypePercentage)

}
