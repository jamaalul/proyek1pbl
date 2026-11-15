package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/jamaalul/proyek1pbl/config"
)

func runWorker() {
	config.LoadEnv()

	srv := asynq.NewServer(
		config.RedisOpt(),
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("[INFO] Background worker Asynq berjalan...")
		if err := srv.Run(mux); err != nil {
			log.Printf("[ERROR] Worker terhenti: %v\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("[INFO] Menghentikan background worker secara graceful...")
	srv.Shutdown()
}
