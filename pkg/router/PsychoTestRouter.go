package router

import (
	"main/pkg/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRouterPsychoTest(controller *controller.PsychoTestController, route *gin.Engine) {

	route.GET("/type-create", controller.CreateType)
	route.POST("/type-addcount", controller.TypeAddCount)
	route.GET("/type-percentage", controller.CountTypePercentage)

}
