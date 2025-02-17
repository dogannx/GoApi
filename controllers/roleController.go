package controllers

import (
	"VueBlog/database"
	"VueBlog/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AllRoles(c *fiber.Ctx) error {
	var roles []models.Role

	database.DB.Find(&roles)

	return c.JSON(roles)
}

func CreateRole(c *fiber.Ctx) error {
	var roleDto fiber.Map

	if err := c.BodyParser(&roleDto); err != nil {
		return err
	}

	//TODO Önemli
	list := roleDto["permissions"].([]interface{})

	permissions := make([]models.Permission, len(list))

	for i, permissionId := range list {
		id := int(permissionId.(float64))
		permissions[i] = models.Permission{
			ID: uint(id),
		}
	}

	role := models.Role{
		Name:        roleDto["name"].(string),
		Permissions: permissions,
	}
	//TODO ÖNEMLi

	database.DB.Create(&role)
	return c.JSON(role)
}

func GetRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	role := models.Role{
		ID: uint(id),
	}
	database.DB.Preload("Permissions").Find(&role)
	return c.JSON(role)
}

func UpdateRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var roleDto fiber.Map
	/*
		role := models.Role{
			ID: uint(id),
		}
	*/
	if err := c.BodyParser(&roleDto); err != nil {
		return err
	}

	//TODO Önemli
	list := roleDto["permissions"].([]interface{})

	permissions := make([]models.Permission, len(list))

	for i, permissionId := range list {
		id := int(permissionId.(float64))
		permissions[i] = models.Permission{
			ID: uint(id),
		}
	}

	type RolePermission struct {
		RoleID       uint
		PermissionID uint
	}

	database.DB.Table("role_permissions").Where("role_id = ?", id).Delete(&RolePermission{})

	role := models.Role{
		ID:          uint(id),
		Name:        roleDto["name"].(string),
		Permissions: permissions,
	}
	//TODO ÖNEMLi

	database.DB.Model(&role).Updates(&role)
	return c.JSON(role)
}

func DeleteRole(c *fiber.Ctx) error {
	// ID'yi al ve dönüşüm hatasını kontrol et
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// Silinecek rolü bul
	role := models.Role{}
	result := database.DB.First(&role, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	// Silme işlemi
	if err := database.DB.Delete(&role).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete roleeeeee",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role deleted successfully",
	})
}
