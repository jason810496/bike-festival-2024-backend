package main

import (
	"main/pkg/config"
	"main/pkg/controller"
	"main/pkg/model"
	"main/pkg/router"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.DatabaseConnection()
	db.Table("phycho_tests").AutoMigrate(&model.PhychoTest{})
	db.Table("calender").AutoMigrate(&model.Calender{})

	// init controller
	phycho_controller := controller.NewPhychoTestController(db)

	// init router
	route := gin.Default()
	router.RegisterRouter_PhychoTest(phycho_controller, route)

	route.Run(":5000")
}
