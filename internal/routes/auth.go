package routes

import (
	"arama-kontrol/internal/handlers"
	"arama-kontrol/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func CreateAuthRoutes(router fiber.Router) {

	router.Post("/signup", handlers.Register)

	router.Post("/login", handlers.Login)

	router.Get("/check-auth", middlewares.VerifyAuth, handlers.CheckAuth)

	router.Get("/logout", handlers.Logout)

	router.Patch("/change-password", middlewares.VerifyAuth, handlers.UpdatePassword)
}
