package landing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mongmx/fiber-cms/models"
	"gorm.io/gorm"
)

func Router(app *fiber.App, db *gorm.DB) {
	repo := NewRepository(db)

	// GET /api/landing-page
	app.Get("/api/landing-page", func(c *fiber.Ctx) error {
		landing, err := repo.GetLandingPage()
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Landing page not found"})
		}

		// Decode JSON strings
		var header, filters, footer map[string]interface{}
		json.Unmarshal([]byte(landing.Header), &header)
		json.Unmarshal([]byte(landing.Filters), &filters)
		json.Unmarshal([]byte(landing.Footer), &footer)

		// Kembalikan data lengkap
		return c.JSON(fiber.Map{
			"header":  header,
			"banner1": landing.Banner1,
			"banner2": landing.Banner2,
			"filters": filters,
			"footer":  footer,
		})
	})

	app.Put("/api/landing-page", func(c *fiber.Ctx) error {
		repo := NewRepository(db)
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form data"})
		}

		// Handle upload banner
		banner1 := ""
		banner2 := ""
		if files := form.File["banner1"]; len(files) > 0 {
			file := files[0]
			fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
			filePath := fmt.Sprintf("./uploads/banners/%s", fileName)
			if err := c.SaveFile(file, filePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload banner1"})
			}
			banner1 = fmt.Sprintf("/uploads/banners/%s", fileName)
		}

		if files := form.File["banner2"]; len(files) > 0 {
			file := files[0]
			fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
			filePath := fmt.Sprintf("./uploads/banners/%s", fileName)
			if err := c.SaveFile(file, filePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload banner2"})
			}
			banner2 = fmt.Sprintf("/uploads/banners/%s", fileName)
		}

		// Ambil data dari form
		header := c.FormValue("header")
		filters := c.FormValue("filters")
		footer := c.FormValue("footer")

		// Ambil data landing dari database
		landing, err := repo.GetLandingPage()
		if err != nil {
			// Jika tidak ada data, buat baru
			landing = &models.Landing{}
		}

		// Update data landing
		if banner1 != "" {
			landing.Banner1 = banner1
		}
		if banner2 != "" {
			landing.Banner2 = banner2
		}
		if header != "" {
			landing.Header = header
		}
		if filters != "" {
			landing.Filters = filters
		}
		if footer != "" {
			landing.Footer = footer
		}

		// Simpan ke database
		if err := repo.UpdateLandingPage(landing); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update landing page"})
		}

		return c.JSON(fiber.Map{"message": "Landing page updated successfully"})
	})
}
