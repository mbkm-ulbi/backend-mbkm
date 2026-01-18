package handlers

import (
	"mbkm-go/internal/database"
	"mbkm-go/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// --- Permissions Handlers ---

func GetPermissions(c *fiber.Ctx) error {
	var permissions []models.Permission
	database.DB.Find(&permissions)
	return c.JSON(fiber.Map{"data": permissions})
}

func CreatePermission(c *fiber.Ctx) error {
	permission := new(models.Permission)
	if err := c.BodyParser(permission); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	database.DB.Create(&permission)
	return c.Status(201).JSON(fiber.Map{"data": permission})
}

func UpdatePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	var permission models.Permission
	if err := database.DB.First(&permission, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	updates := new(models.Permission)
	if err := c.BodyParser(updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	database.DB.Model(&permission).Updates(updates)
	return c.JSON(fiber.Map{"data": permission})
}

func DeletePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.Permission{}, id)
	return c.SendStatus(204)
}

// --- Roles Handlers ---

func GetRoles(c *fiber.Ctx) error {
	var roles []models.Role
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit

	database.DB.Preload("Permissions").Offset(offset).Limit(limit).Find(&roles)
	return c.JSON(fiber.Map{"data": roles})
}

func CreateRole(c *fiber.Ctx) error {
	type RoleInput struct {
		Title       string `json:"title"`
		Permissions []uint `json:"permissions"`
	}
	input := new(RoleInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	role := models.Role{Title: input.Title}
	database.DB.Create(&role)

	if len(input.Permissions) > 0 {
		var perms []models.Permission
		database.DB.Find(&perms, input.Permissions)
		database.DB.Model(&role).Association("Permissions").Replace(&perms)
	}

	return c.Status(201).JSON(fiber.Map{"data": role})
}

func GetRoleDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var role models.Role
	if err := database.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(fiber.Map{"data": role})
}

func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}

	type RoleInput struct {
		Title       string `json:"title"`
		Permissions []uint `json:"permissions"`
	}
	input := new(RoleInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Model(&role).Updates(models.Role{Title: input.Title})

	if len(input.Permissions) > 0 {
		var perms []models.Permission
		database.DB.Find(&perms, input.Permissions)
		database.DB.Model(&role).Association("Permissions").Replace(&perms)
	}

	return c.JSON(fiber.Map{"data": role})
}

func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.Role{}, id)
	return c.SendStatus(204)
}

func AssignRole(c *fiber.Ctx) error {
	type AssignInput struct {
		UserID uint `json:"user_id"`
		RoleID uint `json:"role_id"`
	}
	input := new(AssignInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	if err := database.DB.First(&user, input.UserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	roleName := ""
	switch input.RoleID {
	case 1:
		roleName = "superadmin"
	case 2:
		roleName = "student"
	case 3:
		roleName = "cdc"
	case 4:
		roleName = "company"
	case 5:
		roleName = "dosen"
	case 6:
		roleName = "prodi"
	}

	// Update user string role
	if roleName != "" {
		database.DB.Model(&user).Update("role", roleName)
	}

	// Update many-to-many relationship
	var role models.Role
	if err := database.DB.First(&role, input.RoleID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Role not found"})
	}
	database.DB.Model(&user).Association("Roles").Replace(&role)

	return c.JSON(fiber.Map{"status": true, "data": user})
}
