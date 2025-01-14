package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"

	"pre-test-gallery-service/docs"
	"pre-test-gallery-service/internal/config"
	"pre-test-gallery-service/internal/handlers"
	"pre-test-gallery-service/internal/repository"
	"pre-test-gallery-service/internal/routes"
	"pre-test-gallery-service/internal/service"
	"pre-test-gallery-service/pkg/database"
	"pre-test-gallery-service/pkg/utils"
)

func setupMongoDB(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.ConnectMongoDB(cfg.MongoDBURI)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB successfully")
	return client, nil
}

func setupServer(cfg *config.Config) (*routes.Application, error) {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Go Fiber API v1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	docs.UpdateSwaggerHost(cfg.ServerHost, cfg.ServerPort)
	utils.SetupValidator()

	// Setup CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Authorization,Content-Type",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Setup MongoDB
	mongoClient, err := setupMongoDB(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	db := mongoClient.Database(cfg.MongoDBDatabase)
	tagsRepository := repository.NewTagsRepository(db)

	// // Initialize services
	tagsService := service.NewTagsService(tagsRepository)

	// // Initialize handlers
	tagsHandler := handlers.NewTagsHandler(tagsService)

	// Create application instance
	application := &routes.Application{
		App:         app,
		TagsHandler: tagsHandler,
		Config:      cfg,
	}

	// Setup routes
	application.SetupRoutes()

	return application, nil
}

// @title Service Gallery
// @version 1.0
// @description A RESTful API server
// @termsOfService https://github.com/EyeOfCode

// @contact.name API Support
// @contact.email champuplove@gmail.com

// @host ${DOMAIN}
// @BasePath /api/v1
// @schemes http https
// @in header
func main() {
	cfg := config.LoadConfig()

	application, err := setupServer(cfg)
	if err != nil {
		log.Fatal("Failed to setup server:", err)
	}

	// Graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := cfg.ServerHost + ":" + cfg.ServerPort
		log.Printf("Server starting on %s", addr)
		if err := application.App.Listen(addr); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	if err := application.App.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
