package handlers

import (
	"mbkm-go/internal/database"
	"mbkm-go/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// --- User Handlers ---

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit

	var total int64
	database.DB.Model(&models.User{}).Count(&total)

	database.DB.Preload("Roles").Offset(offset).Limit(limit).Find(&users)
	return c.JSON(fiber.Map{
		"data":  users,
		"count": total,
	})
}

func CreateUser(c *fiber.Ctx) error {
	type UserInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Roles    []uint `json:"roles"`
	}
	input := new(UserInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
		Role:     input.Role, // Default 'student' if empty, handled by DB default usually
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	if len(input.Roles) > 0 {
		var roles []models.Role
		database.DB.Find(&roles, input.Roles)
		database.DB.Model(&user).Association("Roles").Replace(&roles)
	}

	return c.Status(201).JSON(fiber.Map{"data": user})
}

func GetUserDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(fiber.Map{"data": user})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	type UserUpdateInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Roles    []uint `json:"roles"`
		// Add other fields as necessary
	}
	input := new(UserUpdateInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updates := models.User{
		Name:  input.Name,
		Email: input.Email,
	}
	if input.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		updates.Password = string(hash)
	}

	database.DB.Model(&user).Updates(updates)

	if len(input.Roles) > 0 {
		var roles []models.Role
		database.DB.Find(&roles, input.Roles)
		database.DB.Model(&user).Association("Roles").Replace(&roles)
	}

	return c.JSON(fiber.Map{"data": user})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.User{}, id)
	return c.SendStatus(204)
}

// --- Specific Filtering Handlers ---

func GetLecturers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit
	status := c.Query("status")
	applyJobId := c.Query("apply_job_id")

	query := database.DB.Model(&models.User{}).Where("role = ?", "dosen").Preload("Roles")

	if status == "Aktif" {
		query = query.Where("status = ?", "Aktif")
	} else if status == "Tidak Aktif" {
		query = query.Where("status = ?", "Tidak Aktif")
	}

	var total int64
	query.Count(&total)
	var users []models.User
	query.Offset(offset).Limit(limit).Find(&users)

	// Logic for 'lecturer_can_approve' based on applyJobId
	// This requires fetching the ApplyJob and checking ResponsibleLecturerID
	// For API response compatibility, strict struct modification might be needed or just returning map
	var responseData []map[string]interface{}

	// Check responsible lecturer if apply_job_id is present
	var responsibleLecturerID uint = 0
	if applyJobId != "" {
		var applyJob models.ApplyJob
		if err := database.DB.First(&applyJob, applyJobId).Error; err == nil && applyJob.ResponsibleLecturerID != nil {
			responsibleLecturerID = *applyJob.ResponsibleLecturerID
		}
	}

	for _, u := range users {
		userMap := fiber.Map{
			"id":     u.ID,
			"name":   u.Name,
			"email":  u.Email,
			"roles":  u.Roles,
			"status": u.Status,
			// ... add other needed fields
		}
		canApprove := false
		if u.ID == responsibleLecturerID {
			canApprove = true
		}
		userMap["lecturer_can_approve"] = canApprove
		responseData = append(responseData, userMap)
	}

	return c.JSON(fiber.Map{
		"data":  responseData,
		"count": total,
	})
}

func GetStudents(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit
	status := c.Query("status")

	query := database.DB.Model(&models.User{}).Where("role = ?", "student").Preload("Roles")

	switch status {
	case "Aktif":
		query = query.Where("status = ?", "Aktif")
	case "Lulus":
		query = query.Where("status = ?", "Tidak Aktif") // Mapping from Laravel
	case "Drop Out":
		query = query.Where("status = ?", "Drop Out / Dikeluarkan")
	case "Mengundurkan Diri":
		query = query.Where("status = ?", "Mengundurkan Diri / Keluar")
	case "Transfer":
		query = query.Where("status = ?", "Transfer")
	}

	var total int64
	query.Count(&total)
	var users []models.User
	query.Offset(offset).Limit(limit).Find(&users)

	return c.JSON(fiber.Map{
		"data":  users,
		"count": total,
	})
}
