package handlers

import (
	"strconv"

	"mbkm-go/database"
	"mbkm-go/internal/models"
	"mbkm-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct{}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

// Index lists all articles with pagination
func (h *ArticleHandler) Index(c *fiber.Ctx) error {
	page := utils.DefaultPage(c.Query("page"))
	perPage := utils.DefaultLimit(c.Query("per_page"))

	var count int64
	database.DB.Model(&models.Article{}).Count(&count)

	var articles []models.Article
	offset := utils.GetSkipNumber(page, perPage)
	database.DB.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&articles)

	return c.JSON(fiber.Map{
		"data":  articles,
		"count": count,
	})
}

// Show returns a single article
func (h *ArticleHandler) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Article not found")
	}

	var article models.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		return utils.NotFoundError(c, "Article not found")
	}

	// Increment views
	database.DB.Model(&article).Update("views", article.Views+1)

	return c.JSON(fiber.Map{
		"data": article,
	})
}

// Store creates a new article
func (h *ArticleHandler) Store(c *fiber.Ctx) error {
	type ArticleRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req ArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	article := models.Article{
		Title:   req.Title,
		Content: utils.StringPtr(req.Content),
	}

	if err := database.DB.Create(&article).Error; err != nil {
		return utils.InternalServerError(c, "Failed to create article")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": article,
	})
}

// Update updates an article
func (h *ArticleHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Article not found")
	}

	var article models.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		return utils.NotFoundError(c, "Article not found")
	}

	type ArticleRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req ArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	database.DB.Model(&article).Updates(updates)

	return c.JSON(fiber.Map{
		"data": article,
	})
}

// Destroy deletes an article
func (h *ArticleHandler) Destroy(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Article not found")
	}

	var article models.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		return utils.NotFoundError(c, "Article not found")
	}

	database.DB.Delete(&article)

	return c.JSON(fiber.Map{
		"message": "Article deleted successfully",
	})
}
