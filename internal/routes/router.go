package routes

import (
	"pre-test-gallery-service/internal/config"
	"pre-test-gallery-service/internal/handlers"
	"pre-test-gallery-service/pkg/middleware"
	"time"

	"github.com/gofiber/fiber/v2"

	_ "pre-test-gallery-service/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Application struct {
	App              *fiber.App
	TagsHandler      *handlers.TagsHandler
	Config           *config.Config
}

func (app *Application) SetupRoutes() {
	// Swagger route
	app.App.Get("/swagger/*", fiberSwagger.WrapHandler)

	// API routes
	v1 := app.App.Group("/api/v1")
	// Rate limit (You can use route by route)
	v1.Use(middleware.RateLimit(100, time.Minute))

	// Auth routes
	tags := v1.Group("/tags")
	tags.Get("/", app.TagsHandler.GetAllTags)
	tags.Post("/", app.TagsHandler.CreateTags)
	tags.Delete("/:id", app.TagsHandler.DeleteTags)
}