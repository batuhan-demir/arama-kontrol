package main

import (
	"arama-kontrol/internal/dal"
	"arama-kontrol/internal/middlewares"
	"arama-kontrol/internal/routes"
	"arama-kontrol/pkg/database"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	database.Init()

	database.DB.AutoMigrate(&dal.User{}, &dal.Call{}, &dal.Number{})

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	api := app.Group("/api")

	api.Use(middlewares.Protected(), middlewares.GetUser)

	userRouter := api.Group("/users")
	numberRouter := api.Group("/numbers")
	callRouter := api.Group("/calls")
	authRouter := api.Group("/auth")

	routes.CreateUserRoutes(userRouter)
	routes.CreateCallRoutes(callRouter)
	routes.CreateAuthRoutes(authRouter)
	routes.CreateNumberRoutes(numberRouter)

	// react app in dist folder
	app.Static("/", "./dist")
	app.Static("/files", "./files")
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./dist/index.html")
	})

	if os.Getenv("ENV") == "production" {
		app.ListenTLS(":443", "/etc/letsencrypt/live/arama-kontrol.bdemir.net/fullchain.pem", "/etc/letsencrypt/live/arama-kontrol.bdemir.net/privkey.pem")
	} else {
		app.Listen(":" + os.Getenv("PORT"))
	}
}
