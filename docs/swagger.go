package docs

import (
	"fmt"
	"os"
)

// UpdateSwaggerHost updates the Swagger host based on configuration
func UpdateSwaggerHost(host, port string) {
	env := os.Getenv("ENV")
	prodHost := os.Getenv("DOMAIN")

	if env == "production" {
		SwaggerInfo.Host = prodHost
		SwaggerInfo.Schemes = []string{"https"}
	}
	// Update the host directly in SwaggerInfo
	SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, port)
}
