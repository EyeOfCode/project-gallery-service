package docs

import (
	"fmt"
	"os"

	"github.com/swaggo/swag"
)

// @title Fiber API
// @version 1.0
// @description API Documentation

// @host ${SERVER_HOST}:${SERVER_PORT}
// @BasePath /api/v1
// @schemes http https

// SwaggerInfo holds exported Swagger Info
var SwaggerInfo = &swag.Spec{
    Version:     "1.0",
    Host:        fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
    BasePath:    "/api/v1",
    Title:       "Fiber API",
    Description: "API Documentation",
}

func init() {
    swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}