package main

import (
	"log"

	"github.com/RugeFX/ruge-chat-app/database"
	"github.com/RugeFX/ruge-chat-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if envErr := godotenv.Load(); envErr != nil {
		panic(envErr)
	}

	database.ConnectDB()

	app := fiber.New()

	app.Use(cors.New())

	routes.SetupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"hello": "world",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
