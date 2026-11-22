package model

import "time"

type Polyclinic struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreatePolyclinicRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description,omitempty"`
}

type Doctor struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	PolyclinicID   int64     `json:"polyclinic_id"`
	PolyclinicName string    `json:"polyclinic_name,omitempty"`
	Specialization string    `json:"specialization,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateDoctorRequest struct {
	FullName       string `json:"full_name" validate:"required,min=3,max=150"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	PolyclinicID   int64  `json:"polyclinic_id" validate:"required,gt=0"`
	Specialization string `json:"specialization,omitempty"`
}

type DoctorSchedule struct {
	ID             int64     `json:"id"`
	DoctorID       int64     `json:"doctor_id"`
	DoctorName     string    `json:"doctor_name,omitempty"`
	PolyclinicID   int64     `json:"polyclinic_id,omitempty"`
	PolyclinicName string    `json:"polyclinic_name,omitempty"`
	PracticeDate   string    `json:"practice_date"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	Quota          int       `json:"quota"`
	BookedCount    int       `json:"booked_count"`
	RemainingQuota int       `json:"remaining_quota"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateScheduleRequest struct {
	DoctorID     int64  `json:"doctor_id" validate:"required,gt=0"`
	PracticeDate string `json:"practice_date" validate:"required"`
	StartTime    string `json:"start_time" validate:"required"`
	EndTime      string `json:"end_time" validate:"required"`
	Quota        int    `json:"quota" validate:"required,gt=0"`
}

type UpdateScheduleRequest struct {
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
	Quota     *int    `json:"quota,omitempty"`
}
