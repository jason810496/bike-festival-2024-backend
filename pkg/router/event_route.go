package router

import (
	"bikefest/pkg/bootstrap"
	"bikefest/pkg/controller"
)

func RegisterRouterEvent(app *bootstrap.Application, controller *controller.EventController) {

	app.Engine.GET("/event", controller.GetAllEvent)
	app.Engine.GET("/event/:id", controller.GetEventByID)
	app.Engine.POST("/event", controller.CreateEvent)
	app.Engine.PUT("/event/:id", controller.UpdateEvent)
	app.Engine.DELETE("/event/:id", controller.DeleteEvent)
}
