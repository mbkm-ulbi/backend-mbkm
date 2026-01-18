package handlers

import (
	"strconv"

	"mbkm-go/database"
	"mbkm-go/internal/middleware"
	"mbkm-go/internal/models"
	"mbkm-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type CompanyHandler struct{}

func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{}
}

// Index lists all companies with pagination
func (h *CompanyHandler) Index(c *fiber.Ctx) error {
	page := utils.DefaultPage(c.Query("page"))
	perPage := utils.DefaultLimit(c.Query("per_page"))

	var count int64
	database.DB.Model(&models.Company{}).Count(&count)

	var companies []models.Company
	offset := utils.GetSkipNumber(page, perPage)
	database.DB.Preload("User").Preload("CreatedBy").
		Offset(offset).Limit(perPage).Order("created_at DESC").Find(&companies)

	return c.JSON(fiber.Map{
		"data":  companies,
		"count": count,
	})
}

// Show returns a single company
func (h *CompanyHandler) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Company not found")
	}

	var company models.Company
	if err := database.DB.Preload("User").Preload("CreatedBy").First(&company, id).Error; err != nil {
		return utils.NotFoundError(c, "Company not found")
	}

	return c.JSON(fiber.Map{
		"data": company,
	})
}

// Store creates a new company
func (h *CompanyHandler) Store(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return utils.UnauthorizedError(c, "")
	}

	type CompanyRequest struct {
		CompanyName               string `json:"company_name"`
		BusinessFields            string `json:"business_fields"`
		CompanySize               string `json:"company_size"`
		CompanyWebsite            string `json:"company_website"`
		CompanyProfileDescription string `json:"company_profile_description"`
		CompanyPhoneNumber        string `json:"company_phone_number"`
		CompanyAddress            string `json:"company_address"`
	}

	var req CompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

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

	if err := database.DB.Create(&company).Error; err != nil {
		return utils.InternalServerError(c, "Failed to create company")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": company,
	})
}

// Update updates a company
func (h *CompanyHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Company not found")
	}

	var company models.Company
	if err := database.DB.First(&company, id).Error; err != nil {
		return utils.NotFoundError(c, "Company not found")
	}

	type CompanyRequest struct {
		CompanyName               string `json:"company_name"`
		BusinessFields            string `json:"business_fields"`
		CompanySize               string `json:"company_size"`
		CompanyWebsite            string `json:"company_website"`
		CompanyProfileDescription string `json:"company_profile_description"`
		CompanyPhoneNumber        string `json:"company_phone_number"`
		CompanyAddress            string `json:"company_address"`
	}

	var req CompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	updates := map[string]interface{}{}
	if req.CompanyName != "" {
		updates["company_name"] = req.CompanyName
	}
	if req.BusinessFields != "" {
		updates["business_fields"] = req.BusinessFields
	}
	if req.CompanySize != "" {
		updates["company_size"] = req.CompanySize
	}
	if req.CompanyWebsite != "" {
		updates["company_website"] = req.CompanyWebsite
	}
	if req.CompanyProfileDescription != "" {
		updates["company_profile_description"] = req.CompanyProfileDescription
	}
	if req.CompanyPhoneNumber != "" {
		updates["company_phone_number"] = req.CompanyPhoneNumber
	}
	if req.CompanyAddress != "" {
		updates["company_address"] = req.CompanyAddress
	}

	database.DB.Model(&company).Updates(updates)

	return c.JSON(fiber.Map{
		"data": company,
	})
}

// Destroy deletes a company
func (h *CompanyHandler) Destroy(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Company not found")
	}

	var company models.Company
	if err := database.DB.First(&company, id).Error; err != nil {
		return utils.NotFoundError(c, "Company not found")
	}

	database.DB.Delete(&company)

	return c.JSON(fiber.Map{
		"message": "Company deleted successfully",
	})
}
