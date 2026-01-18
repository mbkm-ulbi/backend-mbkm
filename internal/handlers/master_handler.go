package handlers

import (
	"mbkm-go/internal/database"
	"mbkm-go/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// --- Fakultas Handlers ---

func GetFakultas(c *fiber.Ctx) error {
	var fakultas []models.Fakultas

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit

	var total int64
	database.DB.Model(&models.Fakultas{}).Count(&total)

	err := database.DB.Offset(offset).Limit(limit).Find(&fakultas).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}

	return c.JSON(fiber.Map{
		"data":  fakultas,
		"count": total,
	})
}

// --- Program Studi Handlers ---

func GetProgramStudi(c *fiber.Ctx) error {
	var prodi []models.ProgramStudi

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit

	query := database.DB.Model(&models.ProgramStudi{}).Preload("Fakultas")

	// Filter by Fakultas
	if fakultasID := c.Query("fakultas_id"); fakultasID != "" {
		query = query.Where("id_unit_parent = ?", fakultasID)
	}

	var total int64
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Find(&prodi).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}

	return c.JSON(fiber.Map{
		"data":  prodi,
		"count": total,
	})
}

// --- Mata Kuliah Handlers ---

func GetMatkul(c *fiber.Ctx) error {
	var matkul []models.MataKuliah

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit

	query := database.DB.Model(&models.MataKuliah{}).Preload("ProgramStudi")

	// Filter by Prodi
	if prodiID := c.Query("prodi_id"); prodiID != "" {
		query = query.Where("id_program_studi = ?", prodiID)
	}

	var total int64
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).Find(&matkul).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}

	return c.JSON(fiber.Map{
		"data":  matkul,
		"count": total,
	})
}

// --- Perusahaan Handlers ---

func GetPerusahaan(c *fiber.Ctx) error {
	var perusahaan []models.Perusahaan

	// Pagination logic can be added if needed, currently mimicking basic index
	query := database.DB.Model(&models.Perusahaan{}).Preload("User")

	var total int64
	query.Count(&total)

	err := query.Find(&perusahaan).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}

	// Used resource wrapper in Laravel, return key 'data' for consistency if needed,
	// but mostly Laravel resource returns direct object or wrapped in data.
	// UsersApiController returns wrapped data. PerusahaanResource typically wraps.
	return c.JSON(fiber.Map{ // Mimicking Resource collection default wrapping
		"data": perusahaan,
	})
}

func CreatePerusahaan(c *fiber.Ctx) error {
	perusahaan := new(models.Perusahaan)
	if err := c.BodyParser(perusahaan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if result := database.DB.Create(&perusahaan); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": perusahaan})
}

func GetPerusahaanDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var perusahaan models.Perusahaan
	result := database.DB.Preload("User").First(&perusahaan, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Perusahaan not found"})
	}
	return c.JSON(fiber.Map{"data": perusahaan})
}

func UpdatePerusahaan(c *fiber.Ctx) error {
	id := c.Params("id")
	var perusahaan models.Perusahaan
	if result := database.DB.First(&perusahaan, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Perusahaan not found"})
	}

	updates := new(models.Perusahaan)
	if err := c.BodyParser(updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Model(&perusahaan).Updates(updates)
	return c.Status(202).JSON(fiber.Map{"data": perusahaan})
}

func DeletePerusahaan(c *fiber.Ctx) error {
	id := c.Params("id")
	var perusahaan models.Perusahaan
	if result := database.DB.First(&perusahaan, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	database.DB.Delete(&perusahaan)
	return c.SendStatus(204)
}
