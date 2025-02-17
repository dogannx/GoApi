package database

import (
	"VueBlog/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=dogancosman password=157700 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	DB = database

	database.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Product{}, models.Order{}, models.OrderItem{})

}
