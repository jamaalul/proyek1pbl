package model

import "time"

type Role string

const (
	RolePasien Role = "pasien"
	RoleDokter Role = "dokter"
	RoleAdmin  Role = "admin"

	// Legacy aliases if needed
	RoleMahasiswa Role = "mahasiswa"
)

type User struct {
	ID           int64     `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthUser struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	JTI    string `json:"jti,omitempty"`
}
