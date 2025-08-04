package handlers

import (
	"arama-kontrol/internal/dal"
	"arama-kontrol/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetMe(c *fiber.Ctx) error {

	id := c.Locals("user").(jwt.MapClaims)["id"]

	var user dal.User
	database.DB.First(&user, "id = ?", id)

	if user == (dal.User{}) {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"data": &fiber.Map{
			"user": user,
		},
	})
}
