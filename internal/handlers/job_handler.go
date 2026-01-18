package handlers

import (
	"strconv"

	"mbkm-go/database"
	"mbkm-go/internal/dto"
	"mbkm-go/internal/middleware"
	"mbkm-go/internal/models"
	"mbkm-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type JobHandler struct{}

func NewJobHandler() *JobHandler {
	return &JobHandler{}
}

// Index lists all jobs with pagination and filtering
func (h *JobHandler) Index(c *fiber.Ctx) error {
	page := utils.DefaultPage(c.Query("page"))
	perPage := utils.DefaultLimit(c.Query("per_page"))
	companyID := c.Query("company_id")
	status := c.Query("status")

	userID := middleware.GetCurrentUserID(c)
	roleIDs := middleware.GetRoleIDs(c)

	query := database.DB.Model(&models.Job{}).Preload("CreatedBy")

	// Filter by role
	if len(roleIDs) > 0 {
		isAdminOrCDC := false
		for _, rid := range roleIDs {
			if rid == 1 || rid == 3 {
				isAdminOrCDC = true
				break
			}
		}

		if isAdminOrCDC {
			query = query.Where("status IN ?", []string{"Perlu Ditinjau", "Tersedia", "Ditolak"})
		} else {
			query = query.Where("status = ? OR created_by_id = ?", "Tersedia", userID)
		}
	} else {
		query = query.Where("status = ?", "Tersedia")
	}

	// Filter by company
	if companyID != "" {
		if cid, err := strconv.ParseUint(companyID, 10, 32); err == nil {
			query = query.Where("company_id = ?", cid)
		}
	}

	// Filter by status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total
	var count int64
	query.Count(&count)

	// Fetch with pagination
	var jobs []models.Job
	offset := utils.GetSkipNumber(page, perPage)
	query.Offset(offset).Limit(perPage).Order("updated_at DESC").Find(&jobs)

	return c.JSON(fiber.Map{
		"data":  jobs,
		"count": count,
	})
}

// Show returns a single job
func (h *JobHandler) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	var job models.Job
	if err := database.DB.Preload("CreatedBy").First(&job, id).Error; err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	return c.JSON(fiber.Map{
		"data": job,
	})
}

// Store creates a new job
func (h *JobHandler) Store(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return utils.UnauthorizedError(c, "")
	}

	var req dto.JobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	roleIDs := middleware.GetRoleIDs(c)
	isAdminOrCDC := false
	for _, rid := range roleIDs {
		if rid == 1 || rid == 3 {
			isAdminOrCDC = true
			break
		}
	}

	// Get company info if not admin/cdc
	var companyID *uint
	if !isAdminOrCDC {
		var company models.Company
		if err := database.DB.Where("user_id = ?", user.ID).First(&company).Error; err == nil {
			req.Company = company.CompanyName
			if company.CompanyAddress != nil {
				req.Location = *company.CompanyAddress
			}
			companyID = &company.ID
		}
	}

	job := models.Job{
		Title:       req.Title,
		Company:     req.Company,
		Location:    req.Location,
		Duration:    utils.StringPtr(req.Duration),
		Description: utils.StringPtr(req.Description),
		Benefits:    utils.StringPtr(req.Benefits),
		JobType:     utils.StringPtr(req.JobType),
		Salary:      utils.StringPtr(req.Salary),
		VacancyType: utils.StringPtr(req.VacancyType),
		MataKuliah:  utils.StringPtr(req.MataKuliah),
		Deadline:    req.Deadline,
		Status:      models.JobStatusPending,
		CreatedByID: user.ID,
		CompanyID:   companyID,
	}

	if err := database.DB.Create(&job).Error; err != nil {
		return utils.InternalServerError(c, "Failed to create job")
	}

	database.DB.Preload("CreatedBy").First(&job, job.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": job,
	})
}

// Update updates a job
func (h *JobHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	var job models.Job
	if err := database.DB.First(&job, id).Error; err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	var req dto.JobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Duration != "" {
		updates["duration"] = req.Duration
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Benefits != "" {
		updates["benefits"] = req.Benefits
	}
	if req.JobType != "" {
		updates["job_type"] = req.JobType
	}
	if req.Salary != "" {
		updates["salary"] = req.Salary
	}
	if req.VacancyType != "" {
		updates["vacancy_type"] = req.VacancyType
	}
	if req.MataKuliah != "" {
		updates["mata_kuliah"] = req.MataKuliah
	}
	if req.Deadline != nil {
		updates["deadline"] = req.Deadline
	}

	database.DB.Model(&job).Updates(updates)
	database.DB.Preload("CreatedBy").First(&job, job.ID)

	return c.JSON(fiber.Map{
		"data": job,
	})
}

// Destroy deletes a job
func (h *JobHandler) Destroy(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	var job models.Job
	if err := database.DB.First(&job, id).Error; err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	database.DB.Delete(&job)

	return c.JSON(fiber.Map{
		"message": "Data berhasil dihapus",
	})
}

// Approve approves a pending job
func (h *JobHandler) Approve(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	var job models.Job
	if err := database.DB.First(&job, id).Error; err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	result := false
	if job.Status == models.JobStatusPending {
		job.Status = models.JobStatusAvailable
		database.DB.Save(&job)
		result = true
	}

	database.DB.Preload("CreatedBy").First(&job, job.ID)

	return c.JSON(fiber.Map{
		"status": result,
		"data":   job,
	})
}

// Reject rejects a pending job
func (h *JobHandler) Reject(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	var job models.Job
	if err := database.DB.First(&job, id).Error; err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	result := false
	if job.Status == models.JobStatusPending {
		job.Status = models.JobStatusRejected
		database.DB.Save(&job)
		result = true
	}

	database.DB.Preload("CreatedBy").First(&job, job.ID)

	return c.JSON(fiber.Map{
		"status": result,
		"data":   job,
	})
}

// Close closes an available job
func (h *JobHandler) Close(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	var job models.Job
	if err := database.DB.First(&job, id).Error; err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	result := false
	if job.Status == models.JobStatusAvailable {
		job.Status = models.JobStatusClosed
		database.DB.Save(&job)
		result = true
	}

	database.DB.Preload("CreatedBy").First(&job, job.ID)

	return c.JSON(fiber.Map{
		"status": result,
		"data":   job,
	})
}

// ListCandidate lists candidates for a job
func (h *JobHandler) ListCandidate(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.NotFoundError(c, "Job not found")
	}

	// Get apply job IDs for this job
	type ApplyJobJob struct {
		ApplyJobID uint
		JobID      uint
	}
	var applyJobJobs []ApplyJobJob
	database.DB.Table("apply_job_job").Where("job_id = ?", id).Find(&applyJobJobs)

	var candidates []map[string]interface{}
	for _, ajj := range applyJobJobs {
		var applyJob models.ApplyJob
		if err := database.DB.Preload("Users").First(&applyJob, ajj.ApplyJobID).Error; err != nil {
			continue
		}

		for _, user := range applyJob.Users {
			candidates = append(candidates, map[string]interface{}{
				"user":      user,
				"apply_job": applyJob,
			})
		}
	}

	return c.JSON(fiber.Map{
		"data": candidates,
	})
}
