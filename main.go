package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/app/repository"
	"github.com/jamaalul/proyek1pbl/app/service"
	"github.com/jamaalul/proyek1pbl/config"
	"github.com/jamaalul/proyek1pbl/database"
	"github.com/jamaalul/proyek1pbl/route"
)

func main() {
	// 1. Muat variabel lingkungan dari file .env
	config.LoadEnv()

	// 2. Hubungkan ke basis data PostgreSQL
	ctx := context.Background()
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Printf("[PERINGATAN] Database PostgreSQL tidak dapat dihubungi (%v). Mengaktifkan in-memory fallback untuk kebutuhan demo.\n", err)
		pool = nil
	} else {
		defer pool.Close()
		log.Println("[INFO] Terhubung ke basis data PostgreSQL")
	}

	// 3. Inisialisasi Repository dan Service (Clean Architecture tanpa Controller)
	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo)

	// 4. Inisialisasi Fiber Web Framework
	app := fiber.New(fiber.Config{
		AppName: "MediCare Backend API v1.0",
	})

	// 5. Daftarkan seluruh rute endpoint
	route.Register(app, route.Dependencies{
		Pool:        pool,
		AuthService: authService,
	})

	// 6. Jalankan HTTP server
	port := config.GetEnv("APP_PORT", "3000")
	fmt.Printf("\n🚀 MediCare Backend berjalan di http://localhost:%s\n", port)
	fmt.Printf("📋 Health check: http://localhost:%s/api/v1/health\n\n", port)

	log.Fatal(app.Listen(":" + port))
}
