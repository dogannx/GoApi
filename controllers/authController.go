package controllers

import (
	"VueBlog/models"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	user := models.User{
		Name: "Doğan",
	}

	user.Surname = "Coşman"

	return c.JSON(user)
}
