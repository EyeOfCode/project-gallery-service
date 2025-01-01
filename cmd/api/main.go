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

	"go-fiber-api/docs"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/handlers"
	"go-fiber-api/internal/repository"
	"go-fiber-api/internal/routes"
	"go-fiber-api/internal/service"
	"go-fiber-api/pkg/database"
	"go-fiber-api/pkg/utils"
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
		AllowCredentials: cfg.ServerState == "production",
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Setup MongoDB
	mongoClient, err := setupMongoDB(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	db := mongoClient.Database(cfg.MongoDBDatabase)
	userRepository := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepository)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Create application instance
	application := &routes.Application{
		App:            app,
		UserHandler:    userHandler,
		Config:         cfg,
	}

	// Setup routes
	application.SetupRoutes()

	return application, nil
}

// @title Example Go Project API
// @version 1.0
// @description A RESTful API server with user authentication and MongoDB integration
// @termsOfService https://mywebideal.work

// @contact.name API Support
// @contact.email champuplove@gmail.com

// @host localhost:8000
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".
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