package main

import (
	"VueBlog/database"
	"VueBlog/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {

	database.Connect()

	// Initialize a new Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	// Start the server on port 8080
	log.Fatal(app.Listen(":8080"))
}
