package handlers

import (
	"mbkm-go/database"
	"mbkm-go/internal/dto"
	"mbkm-go/internal/middleware"
	"mbkm-go/internal/models"
	"mbkm-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Username == "" || req.Password == "" {
		return utils.ValidationError(c, map[string]string{
			"name":     "Name is required",
			"email":    "Email is required",
			"username": "Username is required",
			"password": "Password is required",
		})
	}

	if len(req.Password) < 6 {
		return utils.ValidationError(c, map[string]string{
			"password": "Password must be at least 6 characters",
		})
	}

	// Check if email/username already exists
	var existingUser models.User
	if err := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		errors := map[string]string{}
		if existingUser.Email == req.Email {
			errors["email"] = "Email already exists"
		}
		if existingUser.Username == req.Username {
			errors["username"] = "Username already exists"
		}
		return utils.ValidationError(c, errors)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.InternalServerError(c, "Failed to hash password")
	}

	// Set default role
	role := "student"
	if req.Role != "" {
		role = req.Role
	}

	// Create user
	user := models.User{
		Name:               req.Name,
		Email:              req.Email,
		Username:           req.Username,
		Password:           string(hashedPassword),
		Role:               role,
		PhoneNumber:        utils.StringPtr(req.PhoneNumber),
		Address:            utils.StringPtr(req.Address),
		ProgramStudy:       utils.StringPtr(req.ProgramStudy),
		Faculty:            utils.StringPtr(req.Faculty),
		NIM:                utils.StringPtr(req.NIM),
		SocialMedia:        utils.StringPtr(req.SocialMedia),
		EmergencyContact:   utils.StringPtr(req.EmergencyContact),
		ProfileDescription: utils.StringPtr(req.ProfileDescription),
		Position:           utils.StringPtr(req.Position),
		Semester:           utils.StringPtr(req.Semester),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.InternalServerError(c, "Failed to create user")
	}

	// Attach role
	var roleModel models.Role
	roleID := uint(2) // Default: student
	switch role {
	case "cdc":
		roleID = 3
	case "mitra", "company":
		roleID = 4
	}

	if err := database.DB.First(&roleModel, roleID).Error; err == nil {
		database.DB.Model(&user).Association("Roles").Append(&roleModel)
	}

	// Create company if role is company/mitra
	if role == "cdc" || role == "company" || role == "mitra" {
		company := models.Company{
			CompanyName:               req.CompanyName,
			BusinessFields:            utils.StringPtr(req.BusinessFields),
			CompanySize:               utils.StringPtr(req.CompanySize),
			CompanyWebsite:            utils.StringPtr(req.CompanyWebsite),
			CompanyProfileDescription: utils.StringPtr(req.CompanyProfileDescription),
			CompanyPhoneNumber:        utils.StringPtr(req.CompanyPhoneNumber),
			CompanyAddress:            utils.StringPtr(req.CompanyAddress),
			UserID:                    &user.ID,
			CreatedByID:               &user.ID,
		}
		database.DB.Create(&company)
	}

	// Create team if not provided
	if req.TeamID == 0 {
		team := models.Team{
			OwnerID: user.ID,
			Name:    req.Email,
		}
		if err := database.DB.Create(&team).Error; err == nil {
			user.TeamID = &team.ID
			database.DB.Save(&user)
		}
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		return utils.InternalServerError(c, "Failed to generate token")
	}

	// Load roles for response
	database.DB.Preload("Roles").First(&user, user.ID)

	roleTitle := ""
	if len(user.Roles) > 0 {
		roleTitle = user.Roles[0].Title
	}

	return c.Status(fiber.StatusCreated).JSON(dto.AuthResponse{
		Token: token,
		Role:  roleTitle,
		User:  user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate
	if req.Username == "" || req.Password == "" {
		return utils.ValidationError(c, map[string]string{
			"username": "Username is required",
			"password": "Password is required",
		})
	}

	// Find user
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).Preload("Roles").First(&user).Error; err != nil {
		return utils.ValidationError(c, map[string]string{
			"username": "Username not found",
		})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.UnauthorizedError(c, "Invalid credentials")
	}

	// Generate token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		return utils.InternalServerError(c, "Failed to generate token")
	}

	roleTitle := ""
	if len(user.Roles) > 0 {
		roleTitle = user.Roles[0].Title
	}

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		Token: token,
		Role:  roleTitle,
		User:  user,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In JWT-based auth, logout is typically handled client-side
	// by removing the token. Server-side can implement token blacklisting if needed.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out",
	})
}

// GetProfile returns current user profile
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return utils.UnauthorizedError(c, "User not authenticated")
	}

	// Get role title
	if len(user.Roles) > 0 {
		user.Role = user.Roles[0].Title
	}

	// Get apply job if exists
	var applyJob models.ApplyJob
	var report *models.Report
	var job *models.Job

	if err := database.DB.Where("created_by_id = ?", user.ID).
		Preload("Jobs").
		Order("created_at DESC").
		First(&applyJob).Error; err == nil {

		// Get report
		database.DB.Where("apply_job_id = ?", applyJob.ID).First(&report)

		// Get job
		if len(applyJob.Jobs) > 0 {
			job = &applyJob.Jobs[0]
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.ProfileResponse{
		User:   user,
		Report: report,
		Job:    job,
	})
}
