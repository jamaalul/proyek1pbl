package config

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() (*slog.Logger, error) {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return nil, err
	}

	file := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     14,
	}
	output := io.MultiWriter(os.Stdout, file)

	return slog.New(slog.NewJSONHandler(output, nil)), nil
}
