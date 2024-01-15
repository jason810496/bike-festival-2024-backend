package router

import (
	"bikefest/pkg/bootstrap"
	"bikefest/pkg/controller"
	"bikefest/pkg/middleware"
)

func RegisterUserRoutes(app *bootstrap.Application, controller *controller.UserController) {
	r := app.Engine.Group("/users")
	authMiddleware := middleware.AuthMiddleware(app.Env.JWT.AccessTokenSecret, app.Cache)

	r.GET("/profile", authMiddleware, controller.Profile)
	r.GET("/:user_id", controller.GetUserByID)
	r.POST("/refresh_token", authMiddleware, controller.RefreshToken)
	r.GET("", controller.GetUsers)
	r.POST("/logout", authMiddleware, controller.Logout)
	r.GET("/login/:user_id", controller.FakeLogin)
	r.POST("/register", controller.FakeRegister)
}
