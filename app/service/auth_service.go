package service

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/app/model"
	"github.com/jamaalul/proyek1pbl/app/repository"
	"github.com/jamaalul/proyek1pbl/config"
	"github.com/jamaalul/proyek1pbl/helper"
)

// AuthService menangani logika bisnis autentikasi dan pendaftaran pengguna
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService membuat instance AuthService baru
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register menangani pendaftaran akun baru
func (s *AuthService) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "Format JSON tidak valid", err.Error())
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.FullName = strings.TrimSpace(req.FullName)

	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" {
		return helper.Error(c, fiber.StatusBadRequest, "Username, email, password, dan nama lengkap wajib diisi", nil)
	}

	if len(req.Password) < 6 {
		return helper.Error(c, fiber.StatusBadRequest, "Password minimal 6 karakter", nil)
	}

	// Default role adalah pasien jika tidak ditentukan
	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role != "admin" && role != "dokter" && role != "pasien" {
		role = "pasien"
	}

	// Cek apakah username sudah dipakai
	if _, err := s.userRepo.FindByUsername(c.Context(), req.Username); err == nil {
		return helper.Error(c, fiber.StatusConflict, "Username sudah digunakan, silakan pilih username lain", nil)
	}

	// Cek apakah email sudah dipakai
	if _, err := s.userRepo.FindByEmail(c.Context(), req.Email); err == nil {
		return helper.Error(c, fiber.StatusConflict, "Email sudah terdaftar", nil)
	}

	// Hash password sebelum disimpan
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "Gagal memproses keamanan password", err.Error())
	}

	newUser := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     role,
		IsActive: true,
	}

	createdUser, err := s.userRepo.Create(c.Context(), newUser)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "Gagal mendaftarkan pengguna baru", err.Error())
	}

	response := model.AuthUser{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		FullName: createdUser.FullName,
		Role:     createdUser.Role,
	}

	return helper.Created(c, "Pendaftaran akun berhasil", response)
}

// Login menangani proses masuk pengguna dan pembuatan access token JWT
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "Format JSON tidak valid", err.Error())
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return helper.Error(c, fiber.StatusBadRequest, "Username dan password wajib diisi", nil)
	}

	user, err := s.userRepo.FindByUsername(c.Context(), req.Username)
	if err != nil {
		return helper.Error(c, fiber.StatusUnauthorized, "Username atau password salah", nil)
	}

	// Verifikasi kata sandi
	if !helper.CheckPasswordHash(req.Password, user.Password) {
		return helper.Error(c, fiber.StatusUnauthorized, "Username atau password salah", nil)
	}

	// Buat JWT Token
	secret := config.GetEnv("JWT_SECRET", "supersecretmedicarejwtkey2026ti")
	expHours := config.GetEnvInt("JWT_EXPIRED_HOURS", 24)

	token, err := helper.GenerateToken(user.ID, user.Username, user.Role, secret, expHours)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "Gagal membuat sesi login token", err.Error())
	}

	tokenResp := model.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		User: model.AuthUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}

	return helper.Success(c, fiber.StatusOK, "Login berhasil", tokenResp)
}

// Profile mengambil data profil pengguna yang sedang login dari context JWT
func (s *AuthService) Profile(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*helper.JWTClaims)
	if !ok || claims == nil {
		return helper.Error(c, fiber.StatusUnauthorized, "Sesi login tidak valid", nil)
	}

	user, err := s.userRepo.FindByID(c.Context(), claims.UserID)
	if err != nil {
		return helper.Error(c, fiber.StatusNotFound, "Data pengguna tidak ditemukan", nil)
	}

	response := model.AuthUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}

	return helper.Success(c, fiber.StatusOK, "Profil pengguna berhasil diambil", response)
}
