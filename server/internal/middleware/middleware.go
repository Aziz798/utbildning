package middleware

import (
	"server/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for Authorization header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authentication token"})
		}

		// The token should be prefixed with "Bearer "
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authentication token"})
		}

		tokenString = tokenParts[1]

		// Verify access token
		claims, err := utils.VerifyToken(tokenString, false)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token, refresh token required"})
		}
		c.Locals("email", claims["email"])
		c.Locals("password", claims["password"])

		// Continue to the next middleware or handler
		return c.Next()
	}
}
