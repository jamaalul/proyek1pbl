package model

import "time"

// User merepresentasikan entitas data pengguna di sistem MediCare
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"` // 'pasien', 'dokter', 'admin'
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
