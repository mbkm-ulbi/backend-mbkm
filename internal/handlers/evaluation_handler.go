package handlers

import (
	"fmt"
	"mbkm-go/internal/database"
	"mbkm-go/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// --- Evaluation Handlers ---

func GetEvaluations(c *fiber.Ctx) error {
	var evaluations []models.Evaluation
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit
	status := c.Query("status")

	query := database.DB.Model(&models.Evaluation{}).
		Preload("ApplyJob.Jobs").
		Preload("ApplyJob.Users").
		Preload("CompanyPersonnel").
		Preload("Lecturer").
		Preload("Examiner").
		Preload("Prodi")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	query.Offset(offset).Limit(limit).Find(&evaluations)

	return c.JSON(fiber.Map{
		"data":  evaluations,
		"count": total,
	})
}

func GetEvaluationDetail(c *fiber.Ctx) error {
	id := c.Params("id") // apply_job_id
	var evaluation models.Evaluation

	if err := database.DB.
		Preload("ApplyJob.Jobs").
		Preload("CompanyPersonnel").
		Preload("Lecturer").
		Preload("Examiner").
		Preload("Prodi").
		Where("apply_job_id = ?", id).
		First(&evaluation).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Evaluation not found"})
	}

	// Calculate Final Grade
	var bobot models.BobotNilai
	database.DB.Where("id_program_studi IS NOT NULL").First(&bobot) // Naive fetch first

	result := calculateFinalGrade(evaluation, bobot)

	return c.JSON(fiber.Map{
		"data": evaluation,
		"meta": result, // Include calculated breakdown
	})
}

func UpdateEvaluation(c *fiber.Ctx) error {
	// Handles grading by Company, Lecturer, Examiner, Prodi
	type GradeInput struct {
		ApplyJobID       uint    `json:"apply_job_id"`
		Grade            string  `json:"grade"` // Letter grade
		GradeScore       float64 `json:"grade_score"`
		GradeDescription string  `json:"grade_description"`
		IsExaminer       bool    `json:"is_examiner"`
	}
	input := new(GradeInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var evaluation models.Evaluation
	if err := database.DB.Where("apply_job_id = ?", input.ApplyJobID).First(&evaluation).Error; err != nil {
		// Create if not exists (though Report usually creates it)
		evaluation = models.Evaluation{
			ApplyJobID: input.ApplyJobID,
			Status:     "Draft",
		}
		database.DB.Create(&evaluation)
	}

	// Calculate Grade from Score
	// gradeChar := "E"
	if input.GradeScore >= 85 {
		// gradeChar = "A"
	} else if input.GradeScore >= 70 {
		// gradeChar = "B"
	} else if input.GradeScore >= 50 {
		// gradeChar = "C"
	} else if input.GradeScore >= 40 {
		// gradeChar = "D"
	}

	// Auth User
	userID := 1 // TODO: Real auth ID
	if u := c.Locals("user_id"); u != nil {
		userID, _ = strconv.Atoi(fmt.Sprintf("%v", u))
	}

	// Determine Role
	var user models.User
	database.DB.Preload("Roles").First(&user, userID)
	roleIDs := []uint{}
	for _, r := range user.Roles {
		roleIDs = append(roleIDs, r.ID)
	}

	now := time.Now()
	updates := map[string]interface{}{}
	updates["status"] = "Sudah Dinilai"

	// Role logic
	for _, rid := range roleIDs {
		if rid == 4 { // Company
			updates["company_personnel_id"] = userID
			updates["company_grade"] = input.Grade // or gradeChar
			updates["company_grade_score"] = input.GradeScore
			updates["company_grade_description"] = input.GradeDescription
			updates["company_grade_date"] = &now
		} else if rid == 5 { // Dosen
			if input.IsExaminer {
				updates["examiner_id"] = userID
				updates["examiner_grade"] = input.Grade
				updates["examiner_grade_score"] = input.GradeScore
				updates["examiner_grade_description"] = input.GradeDescription
				updates["examiner_grade_date"] = &now
			} else {
				updates["lecturer_id"] = userID
				updates["lecturer_grade"] = input.Grade
				updates["lecturer_grade_score"] = input.GradeScore
				updates["lecturer_grade_description"] = input.GradeDescription
				updates["lecturer_grade_date"] = &now
			}
		} else if rid == 6 { // Prodi
			updates["prodi_id"] = userID
			updates["prodi_grade"] = input.Grade
			updates["prodi_grade_score"] = input.GradeScore
			updates["prodi_grade_description"] = input.GradeDescription
			updates["prodi_grade_date"] = &now
		}
	}

	database.DB.Model(&evaluation).Updates(updates)

	return c.JSON(fiber.Map{"data": evaluation, "status": true})
}

// Helpers

func calculateFinalGrade(e models.Evaluation, b models.BobotNilai) map[string]interface{} {
	totalBobot := b.BobotNilaiPerusahaan + b.BobotNilaiPembimbing + b.BobotNilaiPenguji
	if totalBobot == 0 {
		return nil
	}

	pPerusahaan := b.BobotNilaiPerusahaan / totalBobot
	pPembimbing := b.BobotNilaiPembimbing / totalBobot
	pPenguji := b.BobotNilaiPenguji / totalBobot

	cScore := e.CompanyGradeScore * pPerusahaan
	lScore := e.LecturerGradeScore * pPembimbing
	eScore := e.ExaminerGradeScore * pPenguji

	totalScore := cScore + lScore + eScore

	hasGrade := e.CompanyGradeScore > 0 && e.LecturerGradeScore > 0 && e.ExaminerGradeScore > 0

	finalGrade := "E"
	if !hasGrade {
		finalGrade = "-"
	} else if totalScore >= 85 {
		finalGrade = "A"
	} else if totalScore >= 70 {
		finalGrade = "B"
	} else if totalScore >= 55 {
		finalGrade = "C"
	} else if totalScore >= 40 {
		finalGrade = "D"
	}

	return map[string]interface{}{
		"company_score":  cScore,
		"lecturer_score": lScore,
		"examiner_score": eScore,
		"total_score":    totalScore,
		"grade":          finalGrade,
	}
}
