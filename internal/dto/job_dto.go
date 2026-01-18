package dto

import "time"

// JobRequest represents job creation/update request
type JobRequest struct {
	Title       string     `json:"title" validate:"required"`
	Company     string     `json:"company"`
	Location    string     `json:"location"`
	Duration    string     `json:"duration,omitempty"`
	Description string     `json:"description,omitempty"`
	Benefits    string     `json:"benefits,omitempty"`
	JobType     string     `json:"job_type,omitempty"`
	Salary      string     `json:"salary,omitempty"`
	VacancyType string     `json:"vacancy_type,omitempty"`
	MataKuliah  string     `json:"mata_kuliah,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`
}

// JobListRequest represents job list query parameters
type JobListRequest struct {
	Page      int    `query:"page"`
	PerPage   int    `query:"per_page"`
	CompanyID uint   `query:"company_id"`
	Status    string `query:"status"`
}
