package routes

import (
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/handlers"
	"go-fiber-api/pkg/middleware"
	"time"

	"github.com/gofiber/fiber/v2"

	_ "go-fiber-api/pkg/docs"

	"github.com/gofiber/swagger"

	"go-fiber-api/pkg/utils"
)

type Application struct {
	App            *fiber.App
	UserHandler    *handlers.UserHandler
	AuthHandler    *utils.AuthHandler
	Config         *config.Config
}

func (app *Application) SetupRoutes() {
	// API version group
	v1 := app.App.Group("/api/v1")

	// authJwt := utils.NewAuthHandler(app.Config.JWTSecretKey, app.Config.JWTExpiresIn)

	// Global rate limit
	v1.Use(middleware.RateLimit(100, time.Minute))

	// Public routes
	public := v1.Group("")
	{
		// User routes with specific rate limit
		user := public.Group("/user")
		user.Use(middleware.RateLimit(20, time.Minute))
		{
			user.Get("/list", app.UserHandler.UserList)
		}
	}

	// Swagger setup
	app.App.Get("/swagger/*", swagger.HandlerDefault)
}