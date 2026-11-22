package model

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const TaskTypeAppointmentConfirmation = "email:appointment_confirmation"

type AppointmentConfirmationPayload struct {
	AppointmentID int64  `json:"appointment_id"`
	PatientEmail  string `json:"patient_email"`
	PatientName   string `json:"patient_name"`
	QueueNumber   int    `json:"queue_number"`
	PracticeDate  string `json:"practice_date"`
}

func NewAppointmentConfirmationTask(payload AppointmentConfirmationPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal appointment confirmation payload: %w", err)
	}
	return asynq.NewTask(TaskTypeAppointmentConfirmation, data, asynq.Queue("default")), nil
}
