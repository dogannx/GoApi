package controllers

import (
	"VueBlog/database"
	"VueBlog/models"
	"encoding/csv"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
)

func AllOrders(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(database.DB, &models.Order{}, page))
}

func Export(c *fiber.Ctx) error {
	filepath := "./csv/orders.csv"

	if err := CreateFile(filepath); err != nil {
		return err
	}

	return c.Download(filepath)
}

func CreateFile(filepath string) error {
	file, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	writer.Write([]string{
		"ID", "Name", "Email", "Product Title", "Price", "Quantity",
	})

	for _, order := range orders {
		data := []string{
			strconv.Itoa(int(order.ID)),
			order.Firstname + " " + order.Lastname,
			order.Email,
			"",
			"",
			"",
		}
		if err := writer.Write(data); err != nil {
			return err
		}

		for _, orderItem := range order.OrderItems {
			data := []string{
				"",
				"",
				"",
				orderItem.ProductTitle,
				strconv.Itoa(int(orderItem.Price)),
				strconv.Itoa(int(orderItem.Quantity)),
			}
			if err := writer.Write(data); err != nil {
				return err
			}
		}
	}

	return nil
}

type Sales struct {
	Date string `json:"date"`
	Sum  string `json:"sum"`
}

func Chart(c *fiber.Ctx) error {
	var sales []Sales

	database.DB.Raw(`
	SELECT TO_CHAR(o.created_at::timestamp, 'YYYY-MM-DD') AS date,
	SUM(oi.price * oi.quantity) AS sum
	FROM orders o
	JOIN order_items oi ON o.id = oi.order_id
	GROUP BY TO_CHAR(o.created_at::timestamp, 'YYYY-MM-DD');
	`).Scan(&sales)

	return c.JSON(sales)
}
