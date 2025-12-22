package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"` // "SELLER" or "BIDDER"
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
	CompanyID  string `json:"company_id"`
	IsVerified bool   `json:"is_verified"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type VerifyOTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
}

type CreateCompanyRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCompanyRequest struct {
	Name string `json:"name"`
}

type CompanyDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LogoURL    string `json:"logo_url"`
	FoundedDate string `json:"founded_date"`
	Area       string `json:"area"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type Toggle2FARequest struct {
	Enable bool `json:"enable"`
}
