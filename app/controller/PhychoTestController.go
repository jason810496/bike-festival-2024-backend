package controller

import (
	"main/app/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PhychoTestController struct {
	db *gorm.DB
}

func NewPhychoTestController(db *gorm.DB) *PhychoTestController {
	return &PhychoTestController{db: db}
}

// create new phychological type
func (controller *PhychoTestController) CreateType(context *gin.Context) {
	new_type := context.Query("type")

	if new_type == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "Missing or invalid parameters",
		})
		return
	}

	record := model.PhychoTest{
		Type:  new_type,
		Count: 0,
	}

	result := controller.db.Create(&record)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": result.Error,
		})
		panic(result.Error)
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Create type successfully",
	})
}

// Add the count of selected phychological type
func (controller *PhychoTestController) Type_AddCount(context *gin.Context) {
	test_type := context.PostForm("type")
	count, _ := strconv.Atoi(context.PostForm("count"))

	if test_type == "" || count == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "Missing or invalid parameters",
		})
		return
	}

	var phycho_type *model.PhychoTest

	controller.db.Where("type = ?", test_type).First(&phycho_type)

	if phycho_type == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"status":  "Failed",
			"message": "Phychological type doesn't exist",
		})
		return
	}

	phycho_type.Count += count
	controller.db.Save(&phycho_type)
	context.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Successfully add the count of the type",
	})
}

// retrieve the percentage of each type
func (controller *PhychoTestController) CountTypePercentage(context *gin.Context) {
	var query_types []*model.PhychoTest

	controller.db.Find(&query_types)

	if len(query_types) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "Failed",
			"message": "No existing phychological test",
		})
		return
	}

	phycho_types := make(map[string]float64, len(query_types))
	sum := 0

	for _, t := range query_types {
		sum += t.Count
	}

	if sum == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"status":  "Failed",
			"message": "No tested data",
		})
	}

	for _, t := range query_types {
		phycho_types[t.Type] = float64(t.Count) / float64(sum) * 100
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"data":   phycho_types,
	})
}
