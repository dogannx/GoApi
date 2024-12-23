package middlewares

import (
	"VueBlog/database"
	"VueBlog/models"
	"VueBlog/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func IsAuthorized(c *fiber.Ctx, page string) error {
	cookie := c.Cookies("jwt")

	ID, err := util.ParseJwt(cookie)

	if err != nil {
		return err
	}

	userID, _ := strconv.Atoi(ID)

	user := models.User{
		ID: uint(userID),
	}
	database.DB.Preload("Role").Find(&user)

	role := models.Role{
		ID: user.RoleID,
	}
	database.DB.Preload("Permissions").Find(&role)

	if c.Method() == "GET" {
		for _, permission := range role.Permissions {
			if permission.Name == "view_"+page || permission.Name == "edit_"+page {
				return nil
			}
		}
	} else {
		for _, permission := range role.Permissions {
			if permission.Name == "edit_"+page {
				return nil
			}
		}
	}

	c.Status(fiber.StatusUnauthorized)
	return errors.New("unauthorized")
}
