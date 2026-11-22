package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jamaalul/proyek1pbl/app/model"
	"github.com/jamaalul/proyek1pbl/app/service"
	"github.com/jamaalul/proyek1pbl/config"
)

func TestClinicService_Validations(t *testing.T) {
	app := config.NewApp()
	svc := service.NewClinicService(nil, nil)

	app.Post("/schedules", svc.CreateSchedule)
	app.Put("/schedules/:id", svc.UpdateSchedule)

	t.Run("CreateSchedule_PastDateRejected", func(t *testing.T) {
		reqBody, _ := json.Marshal(model.CreateScheduleRequest{
			DoctorID:     1,
			PracticeDate: "2020-01-01",
			StartTime:    "09:00",
			EndTime:      "12:00",
			Quota:        10,
		})

		req := httptest.NewRequest(http.MethodPost, "/schedules", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusUnprocessableEntity)
		}
	})

	t.Run("CreateSchedule_EndBeforeStartRejected", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		futureDate := time.Now().In(loc).AddDate(0, 0, 7).Format("2006-01-02")

		reqBody, _ := json.Marshal(model.CreateScheduleRequest{
			DoctorID:     1,
			PracticeDate: futureDate,
			StartTime:    "14:00",
			EndTime:      "10:00",
			Quota:        10,
		})

		req := httptest.NewRequest(http.MethodPost, "/schedules", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusUnprocessableEntity)
		}
	})

	t.Run("CreateSchedule_QuotaZeroRejected", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		futureDate := time.Now().In(loc).AddDate(0, 0, 7).Format("2006-01-02")

		reqBody, _ := json.Marshal(model.CreateScheduleRequest{
			DoctorID:     1,
			PracticeDate: futureDate,
			StartTime:    "09:00",
			EndTime:      "12:00",
			Quota:        0,
		})

		req := httptest.NewRequest(http.MethodPost, "/schedules", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusUnprocessableEntity)
		}
	})

	t.Run("UpdateSchedule_InvalidIDRejected", func(t *testing.T) {
		reqBody, _ := json.Marshal(model.UpdateScheduleRequest{})
		req := httptest.NewRequest(http.MethodPut, "/schedules/invalid-id", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})
}
