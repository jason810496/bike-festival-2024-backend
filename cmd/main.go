package main

import (
	"bikefest/pkg/bootstrap"
	"bikefest/pkg/router"
	"bikefest/pkg/service"
)

func main() {
	// init config
	app := bootstrap.App()

	// init services
	userService := service.NewUserService(app.Conn, app.Cache)
	eventService := service.NewEventService(app.Conn, app.Cache)

	services := &router.Services{
		UserService:  userService,
		EventService: eventService,
	}

	// init routes
	router.RegisterRoutes(app, services)

	// run app
	app.Run()
}
