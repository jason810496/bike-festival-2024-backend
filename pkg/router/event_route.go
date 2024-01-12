package router

import (
	"bikefest/pkg/bootstrap"
	"bikefest/pkg/controller"
)

func RegisterEventRouter(app *bootstrap.Application, controller *controller.EventController) {
	r := app.Engine.Group("/event")
	r.GET("", controller.GetAllEvent)
	r.GET("/:id", controller.GetEventByID)
	r.POST("", controller.CreateEvent)
	r.PUT("/:id", controller.UpdateEvent)
	r.DELETE("/:id", controller.DeleteEvent)
}
