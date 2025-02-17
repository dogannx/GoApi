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

	// Yeni Fiber uygulamasını başlat
	app := fiber.New()

	// Log middleware'i en başta ekleyelim
	app.Use(func(c *fiber.Ctx) error {
		log.Println("Gelen İstek:", c.Method(), c.OriginalURL())
		log.Println("İstek Gövdesi:", string(c.Body()))
		return c.Next()
	})

	// CORS ayarlarını düzenle
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080", // * yerine frontend adresini ekledik
		AllowMethods:     "GET,POST,PUT,DELETE",   // HTTP metotlarına izin ver
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true, // withCredentials kullanıldığı için true olmalı
	}))

	// Rotaları yükle
	routes.Setup(app)

	// Sunucuyu 8081 portunda başlat
	log.Fatal(app.Listen(":8081"))
}
