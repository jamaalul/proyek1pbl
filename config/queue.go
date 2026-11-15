package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
)

// RedisOpt mengembalikan opsi koneksi Redis dari environment.
func RedisOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       GetEnvInt("REDIS_DB", 0),
	}
}

// NewWorkerServer membuat server Asynq tanpa menjalankannya.
func NewWorkerServer(logger *slog.Logger) *asynq.Server {
	return asynq.NewServer(RedisOpt(), asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		Logger: newAsynqLogger(logger),
	})
}

func newAsynqLogger(logger *slog.Logger) asynq.Logger {
	return &asynqSlogLogger{logger: logger}
}

type asynqSlogLogger struct {
	logger *slog.Logger
}

func (l *asynqSlogLogger) Debug(args ...interface{}) { l.logger.Debug(fmt.Sprint(args...)) }
func (l *asynqSlogLogger) Info(args ...interface{})  { l.logger.Info(fmt.Sprint(args...)) }
func (l *asynqSlogLogger) Warn(args ...interface{})  { l.logger.Warn(fmt.Sprint(args...)) }
func (l *asynqSlogLogger) Error(args ...interface{}) { l.logger.Error(fmt.Sprint(args...)) }

// Fatal mengikuti kontrak Asynq: catat pesan lalu akhiri proses dengan status 1.
func (l *asynqSlogLogger) Fatal(args ...interface{}) {
	l.logger.Error("FATAL: " + fmt.Sprint(args...))
	os.Exit(1)
}
