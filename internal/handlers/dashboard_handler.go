package handlers

import (
	"time"

	"mbkm-go/database"
	"mbkm-go/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// ChartDataset represents a dataset for the chart
type ChartDataset struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

// ChartData represents chart data with labels and datasets
type ChartData struct {
	Labels   []string       `json:"labels"`
	Datasets []ChartDataset `json:"datasets"`
}

// LatestData represents the latest data section
type LatestData struct {
	Jobs             []models.Job      `json:"jobs"`
	Companies        []models.Company  `json:"companies"`
	ApplyJobStudents []models.ApplyJob `json:"apply_job_students"`
}

// DashboardOverview represents the dashboard overview response
type DashboardOverview struct {
	TotalCompany     int64      `json:"total_company"`
	TotalJob         int64      `json:"total_job"`
	TotalStudent     int64      `json:"total_student"`
	TotalAktifMagang int64      `json:"total_aktif_magang"`
	ChartData        ChartData  `json:"chart_data"`
	LatestData       LatestData `json:"latest_data"`
}

// Overview returns dashboard overview data
func (h *DashboardHandler) Overview(c *fiber.Ctx) error {
	var totalCompany, totalJob, totalStudent, totalAktifMagang int64

	// Count totals
	database.DB.Model(&models.Company{}).Count(&totalCompany)
	database.DB.Model(&models.Job{}).Count(&totalJob)
	database.DB.Model(&models.User{}).Where("role = ?", "student").Count(&totalStudent)
	database.DB.Model(&models.ApplyJob{}).Where("status = ?", "Aktif").Count(&totalAktifMagang)

	// Get chart data
	chartData := h.getChartData()

	// Get latest data
	latestData := h.getLatestData()

	data := DashboardOverview{
		TotalCompany:     totalCompany,
		TotalJob:         totalJob,
		TotalStudent:     totalStudent,
		TotalAktifMagang: totalAktifMagang,
		ChartData:        chartData,
		LatestData:       latestData,
	}

	return c.JSON(data)
}

func (h *DashboardHandler) getChartData() ChartData {
	labels := []string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
	statusList := []string{"Melamar", "Disetujui", "Aktif", "Selesai", "Ditolak"}

	var datasets []ChartDataset
	for _, status := range statusList {
		datasets = append(datasets, ChartDataset{
			Label: status,
			Data:  h.getChartDataByStatus(status),
		})
	}

	return ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

func (h *DashboardHandler) getChartDataByStatus(status string) []int {
	type MonthCount struct {
		Month int
		Count int
	}

	var results []MonthCount
	currentYear := time.Now().Year()

	database.DB.Model(&models.ApplyJob{}).
		Select("EXTRACT(MONTH FROM created_at) as month, COUNT(id) as count").
		Where("status = ?", status).
		Where("EXTRACT(YEAR FROM created_at) = ?", currentYear).
		Group("EXTRACT(MONTH FROM created_at)").
		Scan(&results)

	// Format data for all 12 months
	data := make([]int, 12)
	for _, r := range results {
		if r.Month >= 1 && r.Month <= 12 {
			data[r.Month-1] = r.Count
		}
	}

	return data
}

func (h *DashboardHandler) getLatestData() LatestData {
	var jobs []models.Job
	var companies []models.Company
	var applyJobs []models.ApplyJob

	// Get latest 5 jobs
	database.DB.Where("status IN ?", []string{"Perlu Ditinjau", "Tersedia", "Ditolak"}).
		Order("created_at DESC").
		Limit(5).
		Find(&jobs)

	// Get latest 5 companies
	database.DB.Order("created_at DESC").
		Limit(5).
		Find(&companies)

	// Get latest 5 apply jobs with users
	database.DB.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email", "nim", "program_study", "faculty")
	}).
		Order("created_at DESC").
		Limit(5).
		Find(&applyJobs)

	return LatestData{
		Jobs:             jobs,
		Companies:        companies,
		ApplyJobStudents: applyJobs,
	}
}
