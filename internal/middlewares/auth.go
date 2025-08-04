package middlewares

import (	
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("SECRET_KEY"))},
		ErrorHandler: jwtError,
		TokenLookup: "cookie:token", // cookieden al
		ContextKey:  "token",
	})
}

func GetUser(c *fiber.Ctx) error {
	token, _ := c.Locals("token").(*jwt.Token)

	if token == nil {
		c.Locals("user", nil)
		c.Next()
		return nil
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	c.Locals("user", claims)
	c.Next()
	return nil

}

func VerifyAuth(c *fiber.Ctx) error {
	if c.Locals("user") == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Next()
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"success": false,
				"message": "Missing or malformed JWT",
			})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{
			"success": false,
			"message": "Invalid or expired JWT",
		})
}
