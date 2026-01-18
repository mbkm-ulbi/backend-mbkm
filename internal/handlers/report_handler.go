package handlers

import (
	"fmt"
	"mbkm-go/internal/database"
	"mbkm-go/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// --- Report Handlers ---

func GetReports(c *fiber.Ctx) error {
	var reports []models.Report
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit
	status := c.Query("status")

	query := database.DB.Model(&models.Report{}).
		Preload("ApplyJob.Jobs").
		Preload("ApplyJob.Users"). // CreatedBy -> Users
		Preload("CompanyChecked").
		Preload("LecturerChecked").
		Preload("ExaminerChecked").
		Preload("ProdiChecked")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	query.Offset(offset).Limit(limit).Find(&reports)

	return c.JSON(fiber.Map{
		"data":  reports,
		"count": total,
	})
}

func CreateReport(c *fiber.Ctx) error {
	// Parse multipart/form-data
	applyJobID, _ := strconv.Atoi(c.FormValue("apply_job_id"))

	// Check existing report from DB first using apply_job_id
	var existingReport models.Report
	if err := database.DB.Where("apply_job_id = ?", applyJobID).First(&existingReport).Error; err == nil {
		// Report exists, maybe return it or update? Laravel code creates if not exists.
		return c.Status(200).JSON(fiber.Map{"message": "Report already exists", "data": existingReport})
	}

	var applyJob models.ApplyJob
	if err := database.DB.Preload("Jobs").First(&applyJob, applyJobID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Apply Job not found"})
	}

	// TX Start
	tx := database.DB.Begin()

	report := models.Report{
		ApplyJobID:    uint(applyJobID),
		ReportJobUser: *applyJob.JobUser, // Assumed available
		Status:        "Draft",
	}
	if applyJob.CreatedAt.IsZero() == false {
		report.StartDate = &applyJob.CreatedAt
	}

	if err := tx.Create(&report).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Create Evaluation stub
	evaluation := models.Evaluation{
		ApplyJobID: uint(applyJobID),
		Status:     "Belum Dinilai",
	}
	if err := tx.Create(&evaluation).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create evaluation stub"})
	}

	// Handle File Upload
	file, err := c.FormFile("file")
	if err == nil {
		// Save file
		filename := fmt.Sprintf("report_%d_%s", report.ID, file.Filename)
		path := fmt.Sprintf("./uploads/%s", filename) // Ensure uploads dir exists
		if err := c.SaveFile(file, path); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
		}
		// Update report file path/url
		// Assuming simple path storage for now. Laravel uses MediaLibrary.
		updateFile := fmt.Sprintf("/uploads/%s", filename)
		tx.Model(&report).Update("file_laporan", updateFile)
	}

	tx.Commit()

	return c.Status(201).JSON(fiber.Map{"data": report})
}

func GetReportDetail(c *fiber.Ctx) error {
	id := c.Params("id") // apply_job_id usually in Laravel route? Laravel: Route::get('reports/{id}', ...). Controller: show($id) -> where('apply_job_id', $id)
	// Be careful with ID vs ApplyJobID. Laravel `show($id)` uses `Where('apply_job_id', $id)`.
	// Let's support looking up by apply_job_id if the param is used that way.

	var report models.Report
	if err := database.DB.Where("apply_job_id = ?", id).
		Preload("ApplyJob.Jobs").
		Preload("ActivityDetails").
		First(&report).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Report not found"})
	}
	return c.JSON(fiber.Map{"data": report})
}

func CheckReport(c *fiber.Ctx) error {
	id := c.Params("id") // apply_job_id
	var report models.Report
	if err := database.DB.Where("apply_job_id = ?", id).First(&report).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Report not found"})
	}

	var applyJob models.ApplyJob
	database.DB.First(&applyJob, id)

	userID := 1 // TODO: Get actual auth user ID
	// Dummy Auth User ID Retrieval
	//In real app: userID := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"] etc.
	// For now we assume we can get it from context or pass it (but usually from token).
	// We'll rely on a helper if available or placeholder.
	// We'll trust the user passed it in body for testing or fetch from middleware.
	// Let's try to get from Locals if set by middleware, else fail safe.
	if u := c.Locals("user_id"); u != nil {
		userID, _ = strconv.Atoi(fmt.Sprintf("%v", u))
	}

	// Get User Roles
	var user models.User
	database.DB.Preload("Roles").First(&user, userID)

	roleIDs := []uint{}
	for _, r := range user.Roles {
		roleIDs = append(roleIDs, r.ID)
	}

	isExaminer := false
	if applyJob.ExaminerLecturerID != nil && *applyJob.ExaminerLecturerID == uint(userID) {
		isExaminer = true
	}

	now := time.Now()
	updates := map[string]interface{}{}

	// Role 4: Company, 5: Dosen, 6: Prodi
	// Check logic matching Laravel
	for _, rid := range roleIDs {
		if rid == 4 { // Company
			updates["company_checked_id"] = userID
			updates["company_checked_at"] = &now
		} else if rid == 5 { // Dosen
			if isExaminer {
				updates["examiner_checked_id"] = userID
				updates["examiner_checked_at"] = &now
			} else {
				updates["lecturer_checked_id"] = userID
				updates["lecturer_checked_at"] = &now
			}
		} else if rid == 6 { // Prodi
			updates["prodi_checked_id"] = userID
			updates["prodi_checked_at"] = &now
		}
	}

	// Update DB first
	database.DB.Model(&report).Updates(updates)

	// Refetch to check completeness
	database.DB.First(&report, report.ID)

	status := "Berjalan"
	if report.CompanyCheckedID != nil && report.LecturerCheckedID != nil && report.ExaminerCheckedID != nil {
		status = "Selesai"
	}
	database.DB.Model(&report).Update("status", status)
	report.Status = status // Update struct for response

	return c.JSON(fiber.Map{"data": report, "status": true})
}

func DeleteReport(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.Report{}, id)
	return c.SendStatus(204)
}
