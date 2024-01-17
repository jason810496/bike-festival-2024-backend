package router

import (
	"bikefest/pkg/bootstrap"
	"bikefest/pkg/controller"
	"bikefest/pkg/middleware"
)

func RegisterEventRouter(app *bootstrap.Application, controller *controller.EventController) {
	r := app.Engine.Group("/events")
	authMiddleware := middleware.AuthMiddleware(app.Env.JWT.AccessTokenSecret, app.Cache)

	r.GET("", controller.GetAllEvent)
	r.GET("/user", authMiddleware, controller.GetUserEvent)
	r.GET("/:id", controller.GetEventByID)
	r.POST("", controller.SubscribeEvent)
	r.PUT("/:id", controller.UpdateEvent)
	r.DELETE("/:event_id", controller.DeleteEvent)
}
