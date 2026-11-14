package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv memuat file .env jika tersedia
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] File .env tidak ditemukan, menggunakan environment variable bawaan")
	}
}

// GetEnv mengambil nilai env berdasarkan key dengan fallback default
func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// GetEnvInt mengambil nilai env bertipe integer dengan fallback default
func GetEnvInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}
