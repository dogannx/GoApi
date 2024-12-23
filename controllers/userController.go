package controllers

import (
	"VueBlog/database"
	"VueBlog/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	/*	limit := 5
		offset := (page - 1) * limit
		var total int64

		var users []models.User

		database.DB.Preload("Role").Offset(offset).Limit(limit).Find(&users)

		database.DB.Model(&models.User{}).Count(&total)

		return c.JSON(fiber.Map{
			"data": users,
			"meta": fiber.Map{
				"total":     total,
				"page":      page,
				"last_page": math.Ceil(float64(int(total) / limit)),
			},
		})*/ //paginate fonksiyonu ile alakalı burayı kullanmayıp şunu yapıyoruz...
	return c.JSON(models.Paginate(database.DB, &models.User{}, page))
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	//TODO register ile burası farklı orada şifreyi kulllanıcıyı biz belirliyoruz ama gerçekte bunu karşşıdan lamka gerek.
	//password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14) //Bunu sürekli kullanmaktansa user.go nun içerisine bir tane fonksiyon oluşturduk onu kullanacaağız.
	//user.Password = password //Bu satırda alttaki satıra dönüşür.
	user.SetPassword("1234")

	database.DB.Create(&user)
	return c.JSON(user)
}

func GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	user := models.User{
		ID: uint(id),
	}
	database.DB.Preload("Role").Find(&user)
	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	user := models.User{
		ID: uint(id),
	}
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	database.DB.Model(&user).Updates(&user)
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	user := models.User{
		ID: uint(id),
	}
	database.DB.Delete(&user)
	return nil
}
