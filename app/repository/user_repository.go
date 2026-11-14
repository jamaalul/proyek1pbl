package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jamaalul/proyek1pbl/app/model"
)

var (
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("username atau email sudah terdaftar")
)

// UserRepository adalah interface kontrak akses data pengguna
type UserRepository interface {
	Create(ctx context.Context, u model.User) (model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	FindByUsername(ctx context.Context, username string) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
}

// userPostgresRepository implementasi UserRepository menggunakan PostgreSQL (pgxpool)
type userPostgresRepository struct {
	pool *pgxpool.Pool
	// Fallback memory cache untuk memastikan demo tetap bisa berjalan jika database offline
	memMu    sync.RWMutex
	memUsers []model.User
	memNextID int
}

// NewUserRepository membuat instance UserRepository baru
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	repo := &userPostgresRepository{
		pool:      pool,
		memUsers:  make([]model.User, 0),
		memNextID: 1,
	}

	// Seed data demo awal di memory jika DB offline
	if pool == nil {
		repo.memUsers = append(repo.memUsers, model.User{
			ID:        1,
			Username:  "admin",
			Email:     "admin@medicare.local",
			Password:  "$2a$10$7Z8KqK8N49w1g9vS8zG5/O4r5h.g8b9V.U4z5K0L6N1B2V3C4D5E6", // hashed "admin123"
			FullName:  "Administrator Klinik",
			Phone:     "081234567890",
			Role:      "admin",
			IsActive:  true,
			CreatedAt: time.Now(),
		})
		repo.memNextID = 2
	}

	return repo
}

func (r *userPostgresRepository) Create(ctx context.Context, u model.User) (model.User, error) {
	if r.pool != nil {
		err := r.pool.QueryRow(ctx,
			`INSERT INTO users (username, email, password, full_name, phone, role, is_active, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 RETURNING id, created_at`,
			u.Username, u.Email, u.Password, u.FullName, u.Phone, u.Role, u.IsActive, time.Now(),
		).Scan(&u.ID, &u.CreatedAt)

		if err != nil {
			return model.User{}, fmt.Errorf("gagal menyimpan data user ke database: %w", err)
		}
		return u, nil
	}

	// Fallback in-memory
	r.memMu.Lock()
	defer r.memMu.Unlock()

	for _, existing := range r.memUsers {
		if existing.Username == u.Username || existing.Email == u.Email {
			return model.User{}, ErrDuplicate
		}
	}

	u.ID = r.memNextID
	r.memNextID++
	u.CreatedAt = time.Now()
	u.IsActive = true
	r.memUsers = append(r.memUsers, u)
	return u, nil
}

func (r *userPostgresRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	if r.pool != nil {
		var u model.User
		err := r.pool.QueryRow(ctx,
			`SELECT id, username, email, password, full_name, phone, role, is_active, created_at
			 FROM users WHERE id = $1`, id,
		).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.FullName, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.User{}, ErrNotFound
			}
			return model.User{}, err
		}
		return u, nil
	}

	// Fallback in-memory
	r.memMu.RLock()
	defer r.memMu.RUnlock()

	for _, u := range r.memUsers {
		if u.ID == id {
			return u, nil
		}
	}
	return model.User{}, ErrNotFound
}

func (r *userPostgresRepository) FindByUsername(ctx context.Context, username string) (model.User, error) {
	if r.pool != nil {
		var u model.User
		err := r.pool.QueryRow(ctx,
			`SELECT id, username, email, password, full_name, phone, role, is_active, created_at
			 FROM users WHERE LOWER(username) = LOWER($1)`, username,
		).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.FullName, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.User{}, ErrNotFound
			}
			return model.User{}, err
		}
		return u, nil
	}

	// Fallback in-memory
	r.memMu.RLock()
	defer r.memMu.RUnlock()

	for _, u := range r.memUsers {
		if u.Username == username {
			return u, nil
		}
	}
	return model.User{}, ErrNotFound
}

func (r *userPostgresRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	if r.pool != nil {
		var u model.User
		err := r.pool.QueryRow(ctx,
			`SELECT id, username, email, password, full_name, phone, role, is_active, created_at
			 FROM users WHERE LOWER(email) = LOWER($1)`, email,
		).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.FullName, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.User{}, ErrNotFound
			}
			return model.User{}, err
		}
		return u, nil
	}

	// Fallback in-memory
	r.memMu.RLock()
	defer r.memMu.RUnlock()

	for _, u := range r.memUsers {
		if u.Email == email {
			return u, nil
		}
	}
	return model.User{}, ErrNotFound
}
