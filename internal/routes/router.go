package routes

import (
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/handlers"
	"go-fiber-api/pkg/middleware"
	"time"

	"github.com/gofiber/fiber/v2"

	_ "go-fiber-api/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"go-fiber-api/pkg/utils"
)

type Application struct {
	App            *fiber.App
	UserHandler    *handlers.UserHandler
	AuthHandler    *utils.AuthHandler
	ShopHandler    *handlers.ShopHandler
	Config         *config.Config
}

func (app *Application) SetupRoutes() {
	// API version group
	v1 := app.App.Group("/api/v1")

	authJwt := utils.NewAuthHandler(app.Config.JWTSecretKey, app.Config.JWTExpiresIn)

	// Global rate limit
	v1.Use(middleware.RateLimit(100, time.Minute))

	// Public routes
	public := v1.Group("")
	{
		// Auth routes with specific rate limit
		auth := public.Group("/auth")
		auth.Use(middleware.RateLimit(20, time.Minute))
		{
			auth.Post("/login", app.UserHandler.Login)
			auth.Post("/register", app.UserHandler.Register)
		}
	}

	protected := v1.Group("")
	// User routes with jwt
	protected.Use(middleware.JWT(authJwt))
	{
		user := protected.Group("/user")
		{
			user.Put("/profile/:id", app.UserHandler.UpdateProfile)
			user.Get("/profile", app.UserHandler.GetProfile)
			user.Delete("/profile/:id", app.UserHandler.DeleteUser)

			user.Get("/list", app.UserHandler.UserList)
		}

		shop := protected.Group("/shop")
		{
			shop.Get("/list", app.ShopHandler.ShopList)
			shop.Get("/:id", app.ShopHandler.GetShop)
			shop.Post("/", app.ShopHandler.CreateShop)
			shop.Put("/:id", app.ShopHandler.UpdateShop)
			shop.Delete("/:id", app.ShopHandler.DeleteShop)
		}
	}

	// Swagger setup
	app.App.Get("/swagger/*", fiberSwagger.WrapHandler)
}