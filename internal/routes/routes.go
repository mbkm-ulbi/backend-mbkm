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
	dashboardHandler := handlers.NewDashboardHandler()
	applyJobHandler := handlers.NewApplyJobHandler()

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

	// Companies (Job Providers - Complex Entity)
	protectedCompanies := protected.Group("/companies")
	protectedCompanies.Get("", companyHandler.Index)
	protectedCompanies.Get("/:id", companyHandler.Show)
	protectedCompanies.Post("", companyHandler.Store)
	protectedCompanies.Put("/:id", companyHandler.Update)
	protectedCompanies.Delete("/:id", companyHandler.Destroy)

	// --- Master Data Routes ---

	// Fakultas
	protected.Get("/fakultas", handlers.GetFakultas)

	// Program Studi
	protected.Get("/program-studi", handlers.GetProgramStudi)

	// Mata Kuliah
	protected.Get("/matkul", handlers.GetMatkul)

	// Perusahaan (Simple Master Data)
	protectedPerusahaan := protected.Group("/perusahaans")
	protectedPerusahaan.Get("", handlers.GetPerusahaan)
	protectedPerusahaan.Post("", handlers.CreatePerusahaan)
	protectedPerusahaan.Get("/:id", handlers.GetPerusahaanDetail)
	protectedPerusahaan.Put("/:id", handlers.UpdatePerusahaan)
	protectedPerusahaan.Delete("/:id", handlers.DeletePerusahaan)

	// Apply Jobs
	protectedApplyJobs := protected.Group("/apply-jobs")
	protectedApplyJobs.Get("", applyJobHandler.Index)
	protectedApplyJobs.Get("/:id", applyJobHandler.Show)
	protectedApplyJobs.Post("", applyJobHandler.Store)
	protectedApplyJobs.Put("/:id", applyJobHandler.Update)
	protectedApplyJobs.Delete("/:id", applyJobHandler.Destroy)
	protectedApplyJobs.Post("/:id/approve", applyJobHandler.Approve)
	protectedApplyJobs.Post("/:id/reject", applyJobHandler.Reject)
	protectedApplyJobs.Post("/:id/activate", applyJobHandler.Activate)
	protectedApplyJobs.Post("/:id/done", applyJobHandler.Done)
	protectedApplyJobs.Post("/:id/set-lecturer", applyJobHandler.SetLecturer)
	protectedApplyJobs.Get("/user/:user_id", applyJobHandler.GetByUser)

	// Dashboard
	protected.Get("/dashboard/overview", dashboardHandler.Overview)

	// Permissions
	protected.Get("/permissions", handlers.GetPermissions)
	protected.Post("/permissions", handlers.CreatePermission)
	protected.Put("/permissions/:id", handlers.UpdatePermission)
	protected.Delete("/permissions/:id", handlers.DeletePermission)

	// Roles
	protected.Get("/roles", handlers.GetRoles)
	protected.Post("/roles", handlers.CreateRole)
	protected.Get("/roles/:id", handlers.GetRoleDetail)
	protected.Put("/roles/:id", handlers.UpdateRole)
	protected.Delete("/roles/:id", handlers.DeleteRole)
	protected.Post("/roles/assign", handlers.AssignRole)

	// Users
	protected.Get("/users", handlers.GetUsers)
	protected.Post("/users", handlers.CreateUser)
	protected.Get("/users/:id", handlers.GetUserDetail)
	protected.Put("/users/:id", handlers.UpdateUser)
	protected.Delete("/users/:id", handlers.DeleteUser)

	// Special User Filters
	protected.Get("/lecturers", handlers.GetLecturers)
	protected.Get("/students", handlers.GetStudents)

	// --- Academic Features ---

	// Reports
	protectedReports := protected.Group("/reports")
	protectedReports.Get("", handlers.GetReports)
	protectedReports.Post("", handlers.CreateReport)
	protectedReports.Get("/:id", handlers.GetReportDetail) // ID is ApplyJobID
	protectedReports.Post("/:id/check", handlers.CheckReport)
	protectedReports.Delete("/:id", handlers.DeleteReport)

	// Activity Details
	protectedActivities := protected.Group("/activity-details")
	protectedActivities.Get("", handlers.GetActivityDetails)
	protectedActivities.Post("", handlers.CreateActivityDetail)
	protectedActivities.Get("/:id", handlers.GetActivityDetail)
	protectedActivities.Put("/:id", handlers.UpdateActivityDetail)
	protectedActivities.Delete("/:id", handlers.DeleteActivityDetail)

	// Evaluations
	protectedEvaluations := protected.Group("/evaluations")
	protectedEvaluations.Get("", handlers.GetEvaluations)
	protectedEvaluations.Post("", handlers.UpdateEvaluation)       // Store/Update Logic combined
	protectedEvaluations.Get("/:id", handlers.GetEvaluationDetail) // ID is ApplyJobID

	// Konversi Nilai
	protectedKonversi := protected.Group("/konversi-nilai")
	protectedKonversi.Get("", handlers.GetKonversiNilai)
	protectedKonversi.Post("", handlers.CreateKonversiNilai)
	protectedKonversi.Get("/:id", handlers.GetKonversiNilaiDetail)
	protectedKonversi.Put("/:id", handlers.UpdateKonversiNilai)
	protectedKonversi.Delete("/:id", handlers.DeleteKonversiNilai)

	// --- Utilities ---

	// Settings (Bobot Nilai)
	protected.Get("/settings/bobot-nilai", handlers.GetBobotNilai)
	protected.Post("/settings/bobot-nilai", handlers.UpdateBobotNilai)

	// Import
	protected.Post("/import/student", handlers.ImportStudents)
}
