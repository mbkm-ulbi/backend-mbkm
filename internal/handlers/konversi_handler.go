package handlers

import (
	"mbkm-go/internal/database"
	"mbkm-go/internal/models"

	"github.com/gofiber/fiber/v2"
)

// --- Konversi Nilai Handlers ---

func GetKonversiNilai(c *fiber.Ctx) error {
	applyJobID := c.Query("apply_job_id")

	var konversi []models.KonversiNilai
	query := database.DB.Model(&models.KonversiNilai{}).Preload("MataKuliah")

	if applyJobID != "" {
		query = query.Where("apply_job_id = ?", applyJobID)
	}

	if err := query.Find(&konversi).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}

	return c.JSON(fiber.Map{"data": konversi})
}

func CreateKonversiNilai(c *fiber.Ctx) error {
	type KonversiInput struct {
		ApplyJobID uint    `json:"apply_job_id"`
		MatkulID   uint    `json:"matkul_id"`
		Grade      string  `json:"grade"`
		Score      float64 `json:"score"`
	}
	input := new(KonversiInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if exists to update or create (Laravel logic uses find-or-create style in store)
	var konversi models.KonversiNilai
	if err := database.DB.Where("apply_job_id = ? AND matkul_id = ?", input.ApplyJobID, input.MatkulID).First(&konversi).Error; err == nil {
		// Update
		konversi.Grade = input.Grade
		konversi.Score = input.Score
		database.DB.Save(&konversi)
		return c.Status(200).JSON(fiber.Map{"data": konversi, "status": true})
	}

	// Create
	konversi = models.KonversiNilai{
		ApplyJobID: input.ApplyJobID,
		MatkulID:   input.MatkulID,
		Grade:      input.Grade,
		Score:      input.Score,
	}

	if err := database.DB.Create(&konversi).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": konversi, "status": true})
}

func GetKonversiNilaiDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var konversi models.KonversiNilai
	if err := database.DB.Preload("MataKuliah").First(&konversi, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(fiber.Map{"data": konversi})
}

func UpdateKonversiNilai(c *fiber.Ctx) error {
	id := c.Params("id")
	var konversi models.KonversiNilai
	if err := database.DB.First(&konversi, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}

	type UpdateInput struct {
		Grade string  `json:"grade"`
		Score float64 `json:"score"`
	}
	input := new(UpdateInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	konversi.Grade = input.Grade
	konversi.Score = input.Score
	database.DB.Save(&konversi)

	return c.JSON(fiber.Map{"data": konversi})
}

func DeleteKonversiNilai(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.KonversiNilai{}, id)
	return c.SendStatus(204)
}
