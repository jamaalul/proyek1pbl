package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/jamaalul/proyek1pbl/app/model"
	"github.com/jamaalul/proyek1pbl/config"
	"github.com/jamaalul/proyek1pbl/helper"
)

// NewWorkerMux mengonfigurasi multiplexer routing task untuk background worker Asynq.
func NewWorkerMux(mailer helper.Mailer) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.HandleFunc(model.TaskTypeAppointmentConfirmation, func(ctx context.Context, t *asynq.Task) error {
		var payload model.AppointmentConfirmationPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("json unmarshal payload: %w", err)
		}

		subject := fmt.Sprintf("Konfirmasi Reservasi Antrean #%d", payload.QueueNumber)
		body := fmt.Sprintf("Halo %s, reservasi Anda pada tanggal %s berhasil dengan nomor antrean %d.",
			payload.PatientName, payload.PracticeDate, payload.QueueNumber)

		if mailer == nil {
			mailer = helper.NewLogMailer(nil)
		}
		return mailer.Send(ctx, payload.PatientEmail, subject, body)
	})

	return mux
}

func runWorker() {
	config.LoadEnv()

	srv := config.NewWorkerServer(nil)
	mailer := helper.NewLogMailer(nil)
	mux := NewWorkerMux(mailer)

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
