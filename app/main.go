package main

import (
	"main/config"
	"main/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.DatabaseConnection()

	db.Table("phychological_test").AutoMigrate(&model.PhychoTest{})

	route := gin.Default()

	route.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	route.Run(":5000")
}
