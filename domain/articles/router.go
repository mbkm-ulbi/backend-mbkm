package articles

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mongmx/fiber-cms/models"
	"gorm.io/gorm"
)

func Router(app *fiber.App, db *gorm.DB) {
	app.Get("/api/articles", func(c *fiber.Ctx) error { return GetArticles(c, db) })
	app.Get("/api/articles/:id", func(c *fiber.Ctx) error { return GetArticleByID(c, db) })
	app.Post("/api/articles", func(c *fiber.Ctx) error { return CreateArticle(c, db) })
	app.Put("/api/articles/:id", func(c *fiber.Ctx) error { return UpdateArticle(c, db) })
	app.Delete("/api/articles/:id", func(c *fiber.Ctx) error { return DeleteArticle(c, db) })
}

func GetArticles(c *fiber.Ctx, db *gorm.DB) error {
	var articles []models.Article
	if err := db.Find(&articles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	for i := range articles {
		if articles[i].Picture != "" {
			articles[i].Picture = fmt.Sprintf("http://localhost:8080%s", articles[i].Picture)
		}
	}
	return c.JSON(articles)
}

func GetArticleByID(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
	}
	return c.JSON(article)
}

// func CreateArticle(c *fiber.Ctx, db *gorm.DB) error {
// 	article := new(models.Article)
// 	if err := c.BodyParser(article); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
// 	}
// 	if err := db.Create(&article).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	return c.JSON(article)
// }

func CreateArticle(c *fiber.Ctx, db *gorm.DB) error {
	// Parse multipart/form-data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form data"})
	}

	// Ambil data dari form
	title := form.Value["title"][0]
	author := form.Value["author"][0]
	content := form.Value["content"][0]
	// postDate := form.Value["post_date"][0]
	// views := form.Value["views"][0]

	// Handle file upload (gambar)
	var imagePath string
	if fileHeader := form.File["picture"]; len(fileHeader) > 0 {
		file := fileHeader[0] // Ambil file pertama
		fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		fileSavePath := fmt.Sprintf("./uploads/articles/%s", fileName) // Path lokal untuk menyimpan file

		// Simpan file ke folder
		err := c.SaveFile(file, fileSavePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file"})
		}

		// Path relatif untuk menyimpan ke database
		imagePath = fmt.Sprintf("/uploads/articles/%s", fileName)
	}

	// Simpan artikel ke database
	article := models.Article{
		Title:    title,
		Author:   author,
		Content:  content,
		PostDate: time.Now(), // Gunakan `postDate` jika di-parse dari string
		Views:    0,          // Inisialisasi views
		Picture:  imagePath,  // Simpan path gambar
	}
	if err := db.Create(&article).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(article)
}

func UpdateArticle(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
	}
	if err := c.BodyParser(&article); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Save(&article).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(article)
}

// func DeleteArticle(c *fiber.Ctx, db *gorm.DB) error {
// 	id := c.Params("id")
// 	if err := db.Delete(&models.Article{}, id).Error; err != nil {
// 		return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
// 	}
// 	return c.SendStatus(fiber.StatusNoContent)
// }

func DeleteArticle(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
	}

	// Hapus file gambar dari folder
	if article.Picture != "" {
		err := os.Remove("." + article.Picture) // Path lokal file
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete image file"})
		}
	}

	// Hapus artikel dari database
	if err := db.Delete(&article).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete article"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
