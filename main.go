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
	"github.com/jamaalul/proyek1pbl/app/service"
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

	// 3. Setup JWT Manager & Autorisasi RBAC
	jwtSecret := config.GetEnv("JWT_SECRET", "supersecretjwtkey_min32characters_2026")
	jwtIssuer := config.GetEnv("JWT_ISSUER", "4-sehat-5-haji")
	accessTTL := time.Duration(config.GetEnvInt("JWT_TTL_MINUTES", 60)) * time.Minute
	jwtManager := helper.NewJWTManager(jwtSecret, jwtIssuer, accessTTL)

	roleRepo := repository.NewRoleRepository(pool)
	permMap, err := roleRepo.LoadPermissions(ctx)
	if err != nil || len(permMap) == 0 {
		permMap = helper.DefaultClinicPermissions()
	}
	permissions := helper.NewPermissionSet(permMap)

	// 4. Inisialisasi Fiber App dengan Centralized Error Handler & Global Middlewares
	app := config.NewApp()

	txMgr := repository.NewTxManager(pool)
	clinicRepo := repository.NewClinicRepository(pool)
	clinicService := service.NewClinicService(clinicRepo, txMgr)

	// 5. Daftarkan Endpoint Router
	route.Register(app, route.Dependencies{
		Pool:          pool,
		JWT:           jwtManager,
		Permissions:   permissions,
		ClinicService: clinicService,
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

	// 7. Jalankan Asynq Background Worker (Decision #6)
	workerServer := config.NewWorkerServer(nil)
	logMailer := helper.NewLogMailer(nil)
	workerMux := NewWorkerMux(logMailer)

	go func() {
		log.Println("[INFO] Asynq worker server berjalan di background...")
		if err := workerServer.Run(workerMux); err != nil {
			log.Printf("[INFO] Asynq worker berhenti: %v\n", err)
		}
	}()

	// 8. Graceful Shutdown
	<-ctx.Done()
	log.Println("[INFO] Menerima sinyal terminasi, mematikan server dan worker secara graceful...")

	// Matikan Asynq worker
	workerServer.Shutdown()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("[ERROR] Gagal mematikan server Fiber secara normal: %v\n", err)
	}

	log.Println("[INFO] Server dan worker berhasil dimatikan dengan aman")
}
