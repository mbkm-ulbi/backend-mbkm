package jobs

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mongmx/fiber-cms/models"
	"gorm.io/gorm"
)

func Router(app *fiber.App, db *gorm.DB) {
	app.Get("/api/jobs", func(c *fiber.Ctx) error { return GetJobs(c, db) })
	app.Get("/api/jobs/:id", func(c *fiber.Ctx) error { return GetJobByID(c, db) })
	app.Post("/api/jobs", func(c *fiber.Ctx) error { return CreateJob(c, db) })
	app.Put("/api/jobs/:id", func(c *fiber.Ctx) error { return UpdateJob(c, db) })
	app.Delete("/api/jobs/:id", func(c *fiber.Ctx) error { return DeleteJob(c, db) })
}

// Handler Functions
func GetJobs(c *fiber.Ctx, db *gorm.DB) error {
	var jobs []models.Job
	result := db.Find(&jobs)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	for i := range jobs {
		if jobs[i].ImageURL != "" {
			jobs[i].ImageURL = fmt.Sprintf("http://localhost:8080/%s", jobs[i].ImageURL)
		}
	}

	return c.JSON(jobs)
}

func GetJobByID(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	var job models.Job
	if err := db.First(&job, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Job not found"})
	}

	// for i := range jobs {
	// 	if jobs[i].ImageURL != "" {
	// 		jobs[i].ImageURL = fmt.Sprintf("http://localhost:8080/%s", jobs[i].ImageURL)
	// 	}
	// }

	return c.JSON(job)
}

// func CreateJob(c *fiber.Ctx, db *gorm.DB) error {
// 	job := new(models.Job)
// 	if err := c.BodyParser(job); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
// 	}
// 	db.Create(&job)
// 	return c.JSON(job)
// }

func CreateJob(c *fiber.Ctx, db *gorm.DB) error {
	// Parse multipart/form-data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form data"})
	}

	// Get fields from the form
	title := form.Value["title"][0]
	company := form.Value["company"][0]
	location := form.Value["location"][0]
	description := form.Value["description"][0]
	requirements := form.Value["requirements"][0]
	benefits := form.Value["benefits"][0]
	jobType := form.Value["job_type"][0]
	salary := form.Value["salary"][0]

	// Handle file upload
	var imagePath string
	if fileHeader := form.File["image"]; len(fileHeader) > 0 {
		file := fileHeader[0] // Ambil file pertama
		fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		filePath := fmt.Sprintf("./uploads/jobs/%s", fileName)

		// Simpan file ke folder lokal
		err := c.SaveFile(file, filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file"})
		}

		// Simpan path file
		// imagePath = filePath
		imagePath = fmt.Sprintf("uploads/jobs/%s", fileName)
	}

	// Simpan data ke database
	job := models.Job{
		Title:        title,
		Company:      company,
		Location:     location,
		Description:  description,
		Requirements: requirements,
		Benefits:     benefits,
		JobType:      jobType,
		Salary:       salary,
		ImageURL:     imagePath,
	}
	if err := db.Create(&job).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save job"})
	}

	return c.JSON(job)
}

func UpdateJob(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	var job models.Job
	if err := db.First(&job, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Job not found"})
	}
	if err := c.BodyParser(&job); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	db.Save(&job)
	return c.JSON(job)
}

func DeleteJob(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	if err := db.Delete(&models.Job{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Job not found"})
	}
	return c.SendStatus(204)
}
