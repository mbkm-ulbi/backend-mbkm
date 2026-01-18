package handlers

import (
	"strconv"

	"mbkm-go/database"
	"mbkm-go/internal/middleware"
	"mbkm-go/internal/models"
	"mbkm-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplyJobHandler struct{}

func NewApplyJobHandler() *ApplyJobHandler {
	return &ApplyJobHandler{}
}

// Index returns list of apply jobs with pagination
func (h *ApplyJobHandler) Index(c *fiber.Ctx) error {
	page := utils.DefaultPage(c.Query("page"))
	limit := utils.DefaultLimit(c.Query("per_page"))
	skip := utils.GetSkipNumber(page, limit)

	status := c.Query("status")
	companyID := c.Query("company_id")

	query := database.DB.Model(&models.ApplyJob{}).
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "nim", "program_study", "faculty")
		}).
		Preload("Jobs").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Preload("ResponsibleLecturer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Preload("ExaminerLecturer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		})

	// Filter by status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by company_id through jobs
	if companyID != "" {
		query = query.Joins("JOIN apply_job_job ON apply_job_job.apply_job_id = apply_jobs.id").
			Joins("JOIN jobs ON jobs.id = apply_job_job.job_id").
			Where("jobs.company_id = ?", companyID)
	}

	var count int64
	query.Count(&count)

	var applyJobs []models.ApplyJob
	query.Order("created_at DESC").Offset(skip).Limit(limit).Find(&applyJobs)

	return c.JSON(fiber.Map{
		"data":  applyJobs,
		"count": count,
	})
}

// Show returns a single apply job by ID
func (h *ApplyJobHandler) Show(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"id": "Invalid ID"})
	}

	var applyJob models.ApplyJob
	result := database.DB.
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "nim", "program_study", "faculty")
		}).
		Preload("Jobs").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Preload("ResponsibleLecturer").
		Preload("ExaminerLecturer").
		First(&applyJob, id)

	if result.Error != nil {
		return utils.NotFoundError(c, "Apply job not found")
	}

	return c.JSON(fiber.Map{
		"data": applyJob,
	})
}

// ApplyJobRequest for JSON body
type ApplyJobRequest struct {
	Jobs  interface{} `json:"jobs"`
	Users interface{} `json:"users"`
}

// Store creates a new apply job
func (h *ApplyJobHandler) Store(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return utils.UnauthorizedError(c, "Unauthorized")
	}

	var jobIDInt int

	// Try JSON body first
	var req ApplyJobRequest
	if err := c.BodyParser(&req); err == nil && req.Jobs != nil {
		// Handle different types for jobs field
		switch v := req.Jobs.(type) {
		case float64:
			jobIDInt = int(v)
		case int:
			jobIDInt = v
		case string:
			parsedID, err := strconv.Atoi(v)
			if err != nil {
				return utils.ValidationError(c, map[string]string{"jobs": "Invalid job ID"})
			}
			jobIDInt = parsedID
		case []interface{}:
			if len(v) > 0 {
				if id, ok := v[0].(float64); ok {
					jobIDInt = int(id)
				}
			}
		default:
			return utils.ValidationError(c, map[string]string{"jobs": "Invalid job ID format"})
		}
	} else {
		// Fallback to form data
		jobID := c.FormValue("jobs")
		if jobID == "" {
			// Check for jobs[] (frontend sends this)
			jobID = c.FormValue("jobs[]")
		}

		if jobID == "" {
			return utils.ValidationError(c, map[string]string{
				"jobs": "Job ID is required",
			})
		}
		parsedID, err := strconv.Atoi(jobID)
		if err != nil {
			return utils.ValidationError(c, map[string]string{"jobs": "Invalid job ID"})
		}
		jobIDInt = parsedID
	}

	if jobIDInt == 0 {
		return utils.ValidationError(c, map[string]string{"jobs": "Job ID is required"})
	}

	// Check if job exists
	var job models.Job
	if err := database.DB.First(&job, jobIDInt).Error; err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	// Create apply job
	jobUserUUID := uuid.New().String()
	status := "Melamar"
	applyJob := models.ApplyJob{
		JobUser:     &jobUserUUID,
		Status:      &status,
		CreatedByID: &user.ID,
	}

	if err := database.DB.Create(&applyJob).Error; err != nil {
		return utils.InternalServerError(c, "Failed to create apply job")
	}

	// Link user to apply job (many-to-many)
	database.DB.Exec("INSERT INTO apply_job_user (apply_job_id, user_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		applyJob.ID, user.ID)

	// Link job to apply job (many-to-many)
	database.DB.Exec("INSERT INTO apply_job_job (apply_job_id, job_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		applyJob.ID, jobIDInt)

	// Reload with relations
	database.DB.
		Preload("Users").
		Preload("Jobs").
		Preload("CreatedBy").
		First(&applyJob, applyJob.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Apply job created successfully",
		"data":    applyJob,
	})
}

// Update updates an existing apply job
func (h *ApplyJobHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"id": "Invalid ID"})
	}

	var applyJob models.ApplyJob
	if err := database.DB.First(&applyJob, id).Error; err != nil {
		return utils.NotFoundError(c, "Apply job not found")
	}

	// Parse update data
	status := c.FormValue("status")
	lecturerID := c.FormValue("responsible_lecturer_id")

	updates := map[string]interface{}{}
	if status != "" {
		updates["status"] = status
	}
	if lecturerID != "" {
		if lid, err := strconv.Atoi(lecturerID); err == nil {
			updates["responsible_lecturer_id"] = lid
		}
	}

	if examinerID := c.FormValue("examiner_lecturer_id"); examinerID != "" {
		if eid, err := strconv.Atoi(examinerID); err == nil {
			updates["examiner_lecturer_id"] = eid
		}
	}

	if len(updates) > 0 {
		database.DB.Model(&applyJob).Updates(updates)
	}

	// Reload
	database.DB.
		Preload("Users").
		Preload("Jobs").
		Preload("CreatedBy").
		Preload("ResponsibleLecturer").
		Preload("ExaminerLecturer").
		First(&applyJob, id)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    applyJob,
	})
}

// Approve approves an application (Melamar -> Disetujui)
func (h *ApplyJobHandler) Approve(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"id": "Invalid ID"})
	}

	var applyJob models.ApplyJob
	if err := database.DB.First(&applyJob, id).Error; err != nil {
		return utils.NotFoundError(c, "Apply job not found")
	}

	if applyJob.Status == nil || *applyJob.Status != "Melamar" {
		return utils.ValidationError(c, map[string]string{"status": "Can only approve applications with status 'Melamar'"})
	}

	database.DB.Model(&applyJob).Update("status", "Disetujui")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Application approved",
		"data":    applyJob,
	})
}

// Reject rejects an application (Melamar -> Ditolak)
func (h *ApplyJobHandler) Reject(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"id": "Invalid ID"})
	}

	var applyJob models.ApplyJob
	if err := database.DB.First(&applyJob, id).Error; err != nil {
		return utils.NotFoundError(c, "Apply job not found")
	}

	if applyJob.Status == nil || *applyJob.Status != "Melamar" {
		return utils.ValidationError(c, map[string]string{"status": "Can only reject applications with status 'Melamar'"})
	}

	database.DB.Model(&applyJob).Update("status", "Ditolak")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Application rejected",
		"data":    applyJob,
	})
}

// Activate activates an application (Disetujui -> Aktif)
func (h *ApplyJobHandler) Activate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"id": "Invalid ID"})
	}

	var applyJob models.ApplyJob
	if err := database.DB.First(&applyJob, id).Error; err != nil {
		return utils.NotFoundError(c, "Apply job not found")
	}

	if applyJob.Status == nil || *applyJob.Status != "Disetujui" {
		return utils.ValidationError(c, map[string]string{"status": "Can only activate applications with status 'Disetujui'"})
	}

	database.DB.Model(&applyJob).Update("status", "Aktif")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Application activated",
		"data":    applyJob,
	})
}

// Done marks an application as done (Aktif -> Selesai)
func (h *ApplyJobHandler) Done(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"id": "Invalid ID"})
	}

	var applyJob models.ApplyJob
	if err := database.DB.First(&applyJob, id).Error; err != nil {
		return utils.NotFoundError(c, "Apply job not found")
	}

	if applyJob.Status == nil || *applyJob.Status != "Aktif" {
		return utils.ValidationError(c, map[string]string{"status": "Can only mark 'Aktif' applications as done"})
	}

	database.DB.Model(&applyJob).Update("status", "Selesai")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Application completed",
		"data":    applyJob,
	})
}

// SetLecturer assigns a responsible lecturer to an application
func (h *ApplyJobHandler) SetLecturer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"id": "Invalid ID"})
	}

	var applyJob models.ApplyJob
	if err := database.DB.First(&applyJob, id).Error; err != nil {
		return utils.NotFoundError(c, "Apply job not found")
	}

	lecturerID := c.FormValue("lecturer_id")

	updates := map[string]interface{}{}
	if lecturerID != "" {
		if lid, err := strconv.Atoi(lecturerID); err == nil {
			updates["responsible_lecturer_id"] = lid
		}
	}

	examinerID := c.FormValue("examiner_id")
	if examinerID != "" {
		if eid, err := strconv.Atoi(examinerID); err == nil {
			updates["examiner_lecturer_id"] = eid
		}
	}

	// Auto activate if status is Disetujui
	if applyJob.Status != nil && *applyJob.Status == "Disetujui" {
		updates["status"] = "Aktif"
	}

	if len(updates) > 0 {
		database.DB.Model(&applyJob).Updates(updates)
	}

	database.DB.
		Preload("Users").
		Preload("Jobs").
		Preload("ResponsibleLecturer").
		Preload("ExaminerLecturer").
		First(&applyJob, id)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    applyJob,
	})
}

// Destroy deletes an apply job
func (h *ApplyJobHandler) Destroy(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"id": "Invalid ID"})
	}

	var applyJob models.ApplyJob
	if err := database.DB.First(&applyJob, id).Error; err != nil {
		return utils.NotFoundError(c, "Apply job not found")
	}

	database.DB.Delete(&applyJob)

	return c.SendStatus(fiber.StatusNoContent)
}

// GetByUser returns apply jobs for a specific user
func (h *ApplyJobHandler) GetByUser(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return utils.ValidationError(c, map[string]string{"user_id": "Invalid user ID"})
	}

	var applyJobIDs []uint
	database.DB.Table("apply_job_user").
		Where("user_id = ?", userID).
		Pluck("apply_job_id", &applyJobIDs)

	var applyJobs []models.ApplyJob
	database.DB.
		Preload("Users").
		Preload("Jobs").
		Preload("CreatedBy").
		Preload("ResponsibleLecturer").
		Preload("ExaminerLecturer").
		Where("id IN ?", applyJobIDs).
		Find(&applyJobs)

	return c.JSON(fiber.Map{
		"data":  applyJobs,
		"count": len(applyJobs),
	})
}
