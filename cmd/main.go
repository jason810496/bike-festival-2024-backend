package main

import (
	"bikefest/docs"
	"bikefest/pkg/bootstrap"
	"bikefest/pkg/router"
	"bikefest/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"net/http"
	"net/http/httputil"
)

func SetUpSwagger(spec *swag.Spec, app *bootstrap.Application) {
	spec.BasePath = "/"
	spec.Host = fmt.Sprintf("%s:%d", "localhost", app.Env.Server.Port)
	spec.Schemes = []string{"http", "https"}
	spec.Title = "NCKU Bike Festival 2024 Official Website Backend API"
	spec.Description = "This is the official backend API for Bike Festival 2024 Official Website"
}

func ReverseProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			// Copy the original request to retain headers and other attributes
			originalURL := *c.Request.URL

			// Set the scheme and host
			req.URL.Scheme = "http" // or "https" if your target is using SSL
			req.URL.Host = c.Request.Host

			// if the suffix is /docs, then we need to change it to /swagger/index.html
			// otherwise, we need to substitute /docs with /swagger
			if originalURL.Path == "/docs/" {
				req.URL.Path = "/swagger/index.html"
			} else {
				req.URL.Path = "/swagger" + originalURL.Path[len("/docs"):]
			}
			fmt.Println(req.URL.Path)
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

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

	// setup swagger
	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization
	SetUpSwagger(docs.SwaggerInfo, app)
	app.Engine.GET("/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerfiles.Handler,
			ginSwagger.DeepLinking(true),
		),
	)
	app.Engine.GET("/docs/*any", ReverseProxy())

	// run app
	app.Run()
}
