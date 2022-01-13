package middleware

import (
	"os"

	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtected() func(*fiber.Ctx) {
	// Create config for JWT authentication middleware.
	config := jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ErrorHandler: jwtError,
	}

	return jwtware.New(config)
}

func jwtError(ctx *fiber.Ctx, err error) {
	// Return status 400 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
