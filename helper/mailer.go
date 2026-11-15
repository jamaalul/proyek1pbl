package helper

import (
	"context"
	"log/slog"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type LogMailer struct {
	logger *slog.Logger
}

// NewLogMailer membuat mock mailer; logger nil memakai slog.Default().
func NewLogMailer(logger *slog.Logger) *LogMailer {
	if logger == nil {
		logger = slog.Default()
	}
	return &LogMailer{logger: logger}
}

// Send mencatat simulasi tanpa memasukkan isi email ke log.
func (m *LogMailer) Send(ctx context.Context, to, subject, body string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.logger.Info("mail_mock_send", "simulated", true)
	return nil
}

var _ Mailer = (*LogMailer)(nil)
