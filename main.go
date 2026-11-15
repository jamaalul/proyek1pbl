package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jamaalul/proyek1pbl/app/repository"
	"github.com/jamaalul/proyek1pbl/config"
	"github.com/jamaalul/proyek1pbl/database"
	"github.com/jamaalul/proyek1pbl/helper"
	"github.com/jamaalul/proyek1pbl/route"
)

func main() {
	roleFlag := flag.String("role", "api", "Tipe instance: 'api' atau 'worker'")
	flag.Parse()

	if *roleFlag == "worker" {
		runWorker()
		return
	}

	// 1. Muat environment
	config.LoadEnv()

	// 2. Hubungkan ke database PostgreSQL
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Printf("[PERINGATAN] PostgreSQL tidak dapat dihubungi (%v). Berjalan dalam mode fallback.\n", err)
		pool = nil
	} else {
		defer pool.Close()
		log.Println("[INFO] Terhubung ke basis data PostgreSQL")
	}

	// 3. Setup JWT Manager & Autorisasi RBAC (Person 3)
	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Println("[PERINGATAN] JWT_SECRET tidak diatur di environment, menggunakan default development")
		jwtSecret = "supersecretjwtkey_min32characters_2026"
	}
	if len(jwtSecret) < 32 {
		log.Printf("[PERINGATAN] Panjang JWT_SECRET (%d karakter) kurang dari rekomendasi 32 karakter\n", len(jwtSecret))
	}
	jwtIssuer := config.GetEnv("JWT_ISSUER", "siakad-mini")
	accessTTL := time.Duration(config.GetEnvInt("JWT_TTL_MINUTES", 60)) * time.Minute
	jwtManager := helper.NewJWTManager(jwtSecret, jwtIssuer, accessTTL)

	roleRepo := repository.NewRoleRepository(pool)
	permMap, err := roleRepo.LoadPermissions(ctx)
	if err != nil {
		permMap = make(map[string][]string)
	}
	permissions := helper.NewPermissionSet(permMap)

	// 4. Inisialisasi Fiber App dengan Centralized Error Handler & Global Middlewares
	app := config.NewApp()

	// 5. Daftarkan Endpoint Router
	route.Register(app, route.Dependencies{
		Pool:        pool,
		JWT:         jwtManager,
		Permissions: permissions,
	})

	// 6. Jalankan Server HTTP
	port := config.GetEnv("APP_PORT", "3000")
	go func() {
		addr := fmt.Sprintf(":%s", port)
		log.Printf("[INFO] Server HTTP berjalan di http://localhost%s\n", addr)
		if err := app.Listen(addr); err != nil {
			log.Printf("[INFO] Server berhenti: %v\n", err)
		}
	}()

	// 7. Graceful Shutdown
	<-ctx.Done()
	log.Println("[INFO] Menerima sinyal terminasi, mematikan server secara graceful...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("[ERROR] Gagal mematikan server Fiber secara normal: %v\n", err)
	}

	log.Println("[INFO] Server berhasil dimatikan dengan aman")
}
