package routes

import (
	//"arama-kontrol/internal/dal"
	"arama-kontrol/internal/handlers"
	"arama-kontrol/internal/middlewares"
	//"arama-kontrol/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func CreateUserRoutes(router fiber.Router) {

	/*router.Get("/", func(c *fiber.Ctx) error {
		var users []dal.User

		database.DB.Find(&users)

		return c.JSON(users)
	})*/

	router.Get("/me", middlewares.VerifyAuth, handlers.GetMe)

	router.Patch("/me", middlewares.VerifyAuth, handlers.UpdateMe)

}
