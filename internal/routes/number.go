package routes

import (
	"arama-kontrol/internal/handlers"
	"arama-kontrol/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func CreateNumberRoutes(router fiber.Router) {
	router.Get("/", middlewares.VerifyAuth, handlers.GetNumbers)
	
	router.Post("/", middlewares.VerifyAuth, handlers.CreateNumber)

	router.Delete("/:number", middlewares.VerifyAuth, handlers.DeleteNumber)
}
