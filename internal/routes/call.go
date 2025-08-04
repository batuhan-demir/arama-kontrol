package routes

import (
	"arama-kontrol/internal/handlers"
	"arama-kontrol/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func CreateCallRoutes(router fiber.Router) {

	router.Get("/", middlewares.VerifyAuth, handlers.GetCalls)

	router.Post("/callback", handlers.CallCallback)

	router.Patch("/:id/status/:newStatus", middlewares.VerifyAuth, handlers.UpdateCallStatus)

	//router.Get("/receiver/:receiverMail", middlewares.Protected(), middlewares.VerifyToken, handlers.GetFilesByReceiver)

}
