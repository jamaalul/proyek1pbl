package database

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jamaalul/proyek1pbl/config"
)

// NewPool membuat connection pool ke basis data PostgreSQL menggunakan pgxpool
func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	databaseName := config.GetEnv("DB_NAME", "medicare_db")
	dsnURL := url.URL{
		Scheme:  "postgres",
		User:    url.UserPassword(config.GetEnv("DB_USER", "postgres"), config.GetEnv("DB_PASSWORD", "postgres")),
		Host:    net.JoinHostPort(config.GetEnv("DB_HOST", "localhost"), config.GetEnv("DB_PORT", "5432")),
		Path:    "/" + databaseName,
		RawPath: "/" + url.PathEscape(databaseName),
	}
	query := dsnURL.Query()
	query.Set("sslmode", config.GetEnv("DB_SSLMODE", "disable"))
	dsnURL.RawQuery = query.Encode()

	cfg, err := pgxpool.ParseConfig(dsnURL.String())
	if err != nil {
		return nil, fmt.Errorf("konfigurasi PostgreSQL tidak valid; periksa variabel DB_*")
	}

	cfg.MinConns = 2
	maxConns := 10
	if value, ok := os.LookupEnv("DB_MAX_CONNS"); ok && value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("DB_MAX_CONNS harus berupa bilangan bulat")
		}
		maxConns = parsed
	}
	if maxConns < int(cfg.MinConns) || maxConns > math.MaxInt32 {
		return nil, fmt.Errorf("DB_MAX_CONNS harus antara %d dan %d", cfg.MinConns, math.MaxInt32)
	}
	cfg.MaxConns = int32(maxConns)
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat pool database: %w", err)
	}

	// Ping untuk memastikan kredensial benar dan database aktif
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	return pool, nil
}
