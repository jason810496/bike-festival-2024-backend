package main

import (
	"main/pkg/config"
	"main/pkg/controller"
	"main/pkg/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// init config
	app := config.App("config.yaml")

	// init controller
	psychoController := controller.NewPsychoTestController(app.Conn)

	// init router
	route := gin.Default()
	router.RegisterRouterPsychoTest(psychoController, route)

	err := route.Run(":5000")
	if err != nil {
		panic(err)
	}
}
