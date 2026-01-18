package handlers

import (
	"mbkm-go/internal/database"
	"mbkm-go/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// --- ActivityDetail Handlers ---

func GetActivityDetails(c *fiber.Ctx) error {
	var activities []models.ActivityDetail
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit
	reportID := c.Query("report_job_id")

	query := database.DB.Model(&models.ActivityDetail{}).Preload("Report") // Ensure Report relation exists if needed in models

	if reportID != "" {
		query = query.Where("report_job_id = ?", reportID)
	}

	var total int64
	query.Count(&total)

	query.Offset(offset).Limit(limit).Find(&activities)

	return c.JSON(fiber.Map{
		"data":  activities,
		"count": total,
	})
}

func CreateActivityDetail(c *fiber.Ctx) error {
	type ActivityInput struct {
		ReportJobID     uint   `json:"report_job_id"`
		Date            string `json:"date"` // Expects YYYY-MM-DD
		ActivityDetails string `json:"activity_details"`
	}
	input := new(ActivityInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Validate Report Exists
	var report models.Report
	if err := database.DB.First(&report, input.ReportJobID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Report not found"})
	}

	parsedDate, _ := time.Parse("2006-01-02", input.Date)
	activity := models.ActivityDetail{
		ReportJobID:     input.ReportJobID,
		Date:            &parsedDate,
		ActivityDetails: input.ActivityDetails,
	}

	if err := database.DB.Create(&activity).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": activity})
}

func UpdateActivityDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var activity models.ActivityDetail
	if err := database.DB.First(&activity, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}

	type ActivityUpdate struct {
		Date            string `json:"date"`
		ActivityDetails string `json:"activity_details"`
	}
	input := new(ActivityUpdate)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updates := map[string]interface{}{}
	if input.Date != "" {
		parsedDate, _ := time.Parse("2006-01-02", input.Date)
		updates["date"] = &parsedDate
	}
	if input.ActivityDetails != "" {
		updates["activity_details"] = input.ActivityDetails
	}

	database.DB.Model(&activity).Updates(updates)

	return c.JSON(fiber.Map{"data": activity})
}

func DeleteActivityDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.ActivityDetail{}, id)
	return c.SendStatus(204)
}

func GetActivityDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var activity models.ActivityDetail
	if err := database.DB.First(&activity, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(fiber.Map{"data": activity})
}
