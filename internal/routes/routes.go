package routes

import (
	"mbkm-go/internal/handlers"
	"mbkm-go/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	jobHandler := handlers.NewJobHandler()
	articleHandler := handlers.NewArticleHandler()
	companyHandler := handlers.NewCompanyHandler()

	// API v1 routes
	api := app.Group("/api/v1")

	// Test endpoint
	api.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "ok",
		})
	})

	// ================================
	// Public Routes (No Auth Required)
	// ================================

	// Auth routes
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// Public job routes
	public := api.Group("/public")
	public.Get("/jobs", jobHandler.Index)
	public.Get("/jobs/:id", jobHandler.Show)

	// Articles (index and show are public)
	api.Get("/articles", articleHandler.Index)
	api.Get("/articles/:id", articleHandler.Show)

	// ================================
	// Protected Routes (Auth Required)
	// ================================
	protected := api.Group("", middleware.JWTAuth())

	// Auth
	protected.Get("/logout", authHandler.Logout)
	protected.Get("/profile", authHandler.GetProfile)

	// Jobs (with optional auth for filtering)
	jobsWithAuth := api.Group("/jobs", middleware.OptionalJWTAuth())
	jobsWithAuth.Get("", jobHandler.Index)
	jobsWithAuth.Get("/:id", jobHandler.Show)

	// Jobs (protected - create, update, delete, approve, reject, close)
	protectedJobs := protected.Group("/jobs")
	protectedJobs.Post("", jobHandler.Store)
	protectedJobs.Put("/:id", jobHandler.Update)
	protectedJobs.Delete("/:id", jobHandler.Destroy)
	protectedJobs.Post("/:id/approve", jobHandler.Approve)
	protectedJobs.Post("/:id/reject", jobHandler.Reject)
	protectedJobs.Post("/:id/close", jobHandler.Close)
	protectedJobs.Get("/:id/list", jobHandler.ListCandidate)

	// Articles (protected - create, update, delete)
	protectedArticles := protected.Group("/articles")
	protectedArticles.Post("", articleHandler.Store)
	protectedArticles.Put("/:id", articleHandler.Update)
	protectedArticles.Delete("/:id", articleHandler.Destroy)

	// Companies
	protectedCompanies := protected.Group("/companies")
	protectedCompanies.Get("", companyHandler.Index)
	protectedCompanies.Get("/:id", companyHandler.Show)
	protectedCompanies.Post("", companyHandler.Store)
	protectedCompanies.Put("/:id", companyHandler.Update)
	protectedCompanies.Delete("/:id", companyHandler.Destroy)

	// Dashboard (placeholder)
	protected.Get("/dashboard/overview", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Dashboard overview - to be implemented",
		})
	})

	// Permissions (placeholder)
	protected.Get("/permissions", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"data": []interface{}{}, "count": 0})
	})

	// Roles (placeholder)
	protected.Get("/roles", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"data": []interface{}{}, "count": 0})
	})

	// Users (placeholder)
	protected.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"data": []interface{}{}, "count": 0})
	})
	protected.Get("/users/lecturer", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"data": []interface{}{}, "count": 0})
	})
	protected.Get("/users/student", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"data": []interface{}{}, "count": 0})
	})
}
