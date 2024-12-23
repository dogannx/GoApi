package controllers

import (
	"VueBlog/database"
	"VueBlog/models"
	"VueBlog/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

func Register(c *fiber.Ctx) error {

	var data map[string]string

	// Hata Kontrol Yöntem 1
	/*
		err := c.BodyParser(&data)
		if err != nil {
		return err
		}
	*/

	// Hata Kontrol Yöntem 2
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// User oluştruma Yöntem 1
	/*
		user := models.User{
			Name: "veli",
		}
		user.Surname = "Demir"
		user.Email = "veli@gmail.com"
	*/

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"Mesagge": "Passwords do not match",
		})
	}

	//password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14) // Burası da user.go dan dolayı gereksiz bir satır oldu.

	// User oluştruma Yöntem 2
	user := models.User{
		Name:    data["name"],
		Surname: data["surname"],
		Email:   data["email"],
		//Password: password, //burası da user.go içerisindeki fonsksiyon ile alltaki halini alacak...	"user.SetPassword(data["password"])" böyle olur.
		RoleID: 1,
	}

	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	//This section find the user
	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	//This section comparison password on the ready user
	if err := user.ComparePassword(data["password"]); /*bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))*/ err != nil { //Aradaki yorum satırı user.go içerisindeki fonkdan dolayı...
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "You Try False Password",
		})
	}

	/*	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 DAY
		})

		token, err := claims.SignedString([]byte("secret"))
	*/
	// Jwt.go içine yazdığımız fonk sayesinde bu halden  alt satırdakia hale geçtik
	token, err := util.GenerateJwt(strconv.Itoa(int(user.ID)))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	//TODO Burada neden ampersant kullandık sor? && bu coookie kullanımının amacı nedir?

	return c.JSON(fiber.Map{
		"message": "Login Success",
	})

}

// Bu yapıyı kullandığımız yerler standartclaimsi içeriyor aslında hiç kulllanmayıp standartclaim yazabilirsin.
type Claims struct {
	jwt.StandardClaims
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	/*
		token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "Yetkisiz",
			})
		}

		claims := token.Claims.(*Claims)
	*/
	//Üstteki yazımdan alltaki hale...
	id, _ := util.ParseJwt(cookie)

	// claims := token.Claims
	// return c.JSON(claims)
	// Sadece Id Değeri dönüyor

	var user models.User
	database.DB.Where("id = ?", id).First(&user)
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "Logout Success",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	cookie := c.Cookies("jwt")
	id, _ := util.ParseJwt(cookie)

	userId, _ := strconv.Atoi(id)

	user := models.User{
		ID:      uint(userId),
		Name:    data["name"],
		Surname: data["surname"],
		Email:   data["email"],
	}
	database.DB.Model(&user).Updates(user)

	return c.JSON(user)

}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"Mesagge": "Passwords do not match",
		})
	}

	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)
	
	userId, _ := strconv.Atoi(id)

	user := models.User{
		ID: uint(userId),
	}

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)

}
