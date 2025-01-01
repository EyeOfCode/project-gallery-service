package middleware

import (
	"go-fiber-api/pkg/utils"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWT(auth *utils.AuthHandler, role ...utils.Role) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return utils.SendError(c, http.StatusUnauthorized, "Authorization header is required")
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
            return utils.SendError(c, http.StatusUnauthorized, "Invalid token format")
        }

        token := bearerToken[1]
        claims, err := auth.ValidateToken(token)
        if err != nil {
            return utils.SendError(c, http.StatusUnauthorized, "Invalid token")
        }

        if len(role) > 0 {
            roleSlice := make([]utils.Role, len(claims.Roles))
            for i, r := range claims.Roles {
                roleSlice[i] = utils.Role(r)
            }

            if !utils.IsValidRole(roleSlice, role) {
                return utils.SendError(c, http.StatusForbidden, "Insufficient permissions")
            }
        }

        c.Locals("userID", claims.UserID)
        return c.Next()
    }
}