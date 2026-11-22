package model_test

import (
	"encoding/json"
	"testing"

	"github.com/jamaalul/proyek1pbl/app/model"
)

func TestNewAppointmentConfirmationTask(t *testing.T) {
	payload := model.AppointmentConfirmationPayload{
		AppointmentID: 101,
		PatientEmail:  "patient@example.com",
		PatientName:   "Budi Santoso",
		QueueNumber:   5,
		PracticeDate:  "2026-12-01",
	}

	task, err := model.NewAppointmentConfirmationTask(payload)
	if err != nil {
		t.Fatalf("unexpected error creating task: %v", err)
	}

	if task.Type() != model.TaskTypeAppointmentConfirmation {
		t.Errorf("task type = %s, want %s", task.Type(), model.TaskTypeAppointmentConfirmation)
	}

	var parsed model.AppointmentConfirmationPayload
	if err := json.Unmarshal(task.Payload(), &parsed); err != nil {
		t.Fatalf("failed to unmarshal task payload: %v", err)
	}

	if parsed.AppointmentID != 101 || parsed.QueueNumber != 5 || parsed.PatientEmail != "patient@example.com" {
		t.Errorf("payload mismatch: %+v", parsed)
	}
}
