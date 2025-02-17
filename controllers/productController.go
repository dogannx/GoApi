package controllers

import (
	"VueBlog/database"
	"VueBlog/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AllProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(database.DB, &models.Product{}, page))
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	// ID sıfırlanmasına gerek yok, GORM bunu otomatik yapar
	result := database.DB.Create(&product)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create product",
			"details": result.Error.Error(),
		})
	}

	return c.Status(201).JSON(product) // 201 status code ürün oluşturulduğunu belirtir
}

func GetProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	product := models.Product{
		ID: uint(id),
	}
	database.DB.Find(&product)
	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	product := models.Product{
		ID: uint(id),
	}
	if err := c.BodyParser(&product); err != nil {
		return err
	}
	database.DB.Model(&product).Updates(&product)
	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	product := models.Product{
		ID: uint(id),
	}
	database.DB.Delete(&product)
	return nil
}
