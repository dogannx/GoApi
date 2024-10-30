package routes

import (
	"VueBlog/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Define a route for the GET method on the root path '/'
	app.Post("/api/register", controllers.Register)
}
