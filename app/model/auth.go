package model

// RegisterRequest adalah DTO data pendaftaran akun baru
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"` // 'pasien', 'dokter', 'admin'
}

// LoginRequest adalah DTO data kredensial login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthUser adalah representasi data pengguna publik yang aman dikembalikan ke klien
type AuthUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// TokenResponse adalah payload respon saat login berhasil
type TokenResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	User        AuthUser `json:"user"`
}
