package dto

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Name               string `json:"name" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password" validate:"required,min=6"`
	PhoneNumber        string `json:"phone_number,omitempty"`
	Address            string `json:"address,omitempty"`
	ProgramStudy       string `json:"program_study,omitempty"`
	Faculty            string `json:"faculty,omitempty"`
	NIM                string `json:"nim,omitempty"`
	Semester           string `json:"semester,omitempty"`
	SocialMedia        string `json:"social_media,omitempty"`
	EmergencyContact   string `json:"emergency_contact,omitempty"`
	ProfileDescription string `json:"profile_description,omitempty"`
	Position           string `json:"position,omitempty"`
	Role               string `json:"role,omitempty"` // student, cdc, company, mitra
	TeamID             uint   `json:"team,omitempty"`
	// Company fields
	CompanyName               string `json:"company_name,omitempty"`
	BusinessFields            string `json:"business_fields,omitempty"`
	CompanySize               string `json:"company_size,omitempty"`
	CompanyWebsite            string `json:"company_website,omitempty"`
	CompanyProfileDescription string `json:"company_profile_description,omitempty"`
	CompanyPhoneNumber        string `json:"company_phone_number,omitempty"`
	CompanyAddress            string `json:"company_address,omitempty"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string      `json:"token"`
	Role  string      `json:"role"`
	User  interface{} `json:"user"`
}

// ProfileResponse represents profile response
type ProfileResponse struct {
	User   interface{} `json:"user"`
	Report interface{} `json:"report"`
	Job    interface{} `json:"job"`
}
