package route_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jamaalul/proyek1pbl/app/service"
	"github.com/jamaalul/proyek1pbl/config"
	"github.com/jamaalul/proyek1pbl/helper"
	"github.com/jamaalul/proyek1pbl/route"
)

func TestClinicRoutes_AuthAndPermissions(t *testing.T) {
	app := config.NewApp()
	jwtMgr := helper.NewJWTManager("testsecretkey_min32characters_12345", "test-issuer", time.Hour)
	perms := helper.NewPermissionSet(helper.DefaultClinicPermissions())
	svc := service.NewClinicService(nil, nil)

	v1 := app.Group("/api/v1")
	route.RegisterClinicRoutes(v1, svc, jwtMgr, perms)

	t.Run("UnauthorizedWithoutToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/polyclinics", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("PatientCannotManageClinic_Forbidden", func(t *testing.T) {
		token, err := jwtMgr.GenerateTokenFor(10, "patient@test.com", "pasien")
		if err != nil {
			t.Fatalf("generate token failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/polyclinics", bytes.NewReader([]byte(`{"name":"Mata"}`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
	})

	t.Run("PatientCannotCreateDoctor_Forbidden", func(t *testing.T) {
		token, err := jwtMgr.GenerateTokenFor(10, "patient@test.com", "pasien")
		if err != nil {
			t.Fatalf("generate token failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/doctors", bytes.NewReader([]byte(`{"full_name":"Dr"}`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
	})

	t.Run("PatientCannotCreateSchedule_Forbidden", func(t *testing.T) {
		token, err := jwtMgr.GenerateTokenFor(10, "patient@test.com", "pasien")
		if err != nil {
			t.Fatalf("generate token failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewReader([]byte(`{"doctor_id":1}`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
	})
}
