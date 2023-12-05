package main

import (
	"main/config"
	"main/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.DatabaseConnection()

	db.Table("phycho_tests").AutoMigrate(&model.PhychoTest{})

	route := gin.Default()

	route.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	route.POST("/add_test", func(context *gin.Context) {
		test_type := context.PostForm("type")
		count, _ := strconv.Atoi(context.PostForm("count"))

		record := model.PhychoTest{
			Type:  test_type,
			Count: count,
		}

		result := db.Create(&record)

		if result.Error != nil {
			context.JSON(400, gin.H{
				"status":  "post failed",
				"type":    test_type,
				"count":   count,
				"message": result.Error,
			})
			panic(result.Error)
		}

		context.JSON(http.StatusOK, gin.H{
			"status": "posted",
			"type":   test_type,
			"count":  count,
		})
	})

	route.GET("/get_type", func(context *gin.Context) {
		test_type := context.Query("type")

		var tests []*model.PhychoTest

		db.Where("type = ?", test_type).Find(&tests)

		context.JSON(http.StatusOK, gin.H{
			"status":        "retrieved",
			"type":          test_type,
			"num_of_result": len(tests),
		})
	})

	route.Run(":5000")
}
