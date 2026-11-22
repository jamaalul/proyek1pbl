package main

import (
	"context"
	"testing"

	"github.com/jamaalul/proyek1pbl/app/model"
)

type mockMailer struct {
	sentTo   string
	sentSub  string
	sentBody string
}

func (m *mockMailer) Send(ctx context.Context, to, subject, body string) error {
	m.sentTo = to
	m.sentSub = subject
	m.sentBody = body
	return nil
}

func TestWorkerMux_AppointmentConfirmation(t *testing.T) {
	mock := &mockMailer{}
	mux := NewWorkerMux(mock)

	payload := model.AppointmentConfirmationPayload{
		AppointmentID: 202,
		PatientEmail:  "patient@test.com",
		PatientName:   "Ahmad",
		QueueNumber:   3,
		PracticeDate:  "2026-12-05",
	}

	task, err := model.NewAppointmentConfirmationTask(payload)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	err = mux.ProcessTask(context.Background(), task)
	if err != nil {
		t.Fatalf("expected nil error processing task, got: %v", err)
	}

	if mock.sentTo != "patient@test.com" {
		t.Errorf("sentTo = %s, want patient@test.com", mock.sentTo)
	}
}
