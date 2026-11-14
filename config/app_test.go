package config_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/app/model"
	"github.com/jamaalul/proyek1pbl/config"
	"github.com/jamaalul/proyek1pbl/helper"
	"github.com/jamaalul/proyek1pbl/middleware"
)

func TestApp_CentralizedErrorHandler(t *testing.T) {
	app := config.NewApp()

	app.Get("/trigger-404", func(c *fiber.Ctx) error {
		return helper.NotFound("data tidak ditemukan")
	})

	app.Get("/trigger-validation", func(c *fiber.Ctx) error {
		return helper.Validation(map[string][]string{
			"email": {"email wajib diisi"},
		})
	})

	app.Get("/trigger-panic", func(c *fiber.Ctx) error {
		panic("simulasi crash server")
	})

	t.Run("HandleAppErrorNotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/trigger-404", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}

		var errRes model.ErrorResponse
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &errRes); err != nil {
			t.Fatalf("gagal unmarshal ErrorResponse: %v", err)
		}
		if errRes.Success || errRes.Message != "data tidak ditemukan" {
			t.Errorf("ErrorResponse tidak sesuai: %+v", errRes)
		}
	})

	t.Run("HandleValidationError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/trigger-validation", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusUnprocessableEntity)
		}

		var errRes model.ErrorResponse
		body, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &errRes)
		if errRes.Errors["email"][0] != "email wajib diisi" {
			t.Errorf("Validation errors map tidak sesuai: %+v", errRes.Errors)
		}
	})

	t.Run("HandlePanicRecovery", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/trigger-panic", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})

	t.Run("RequestIDHeaderPresence", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/trigger-404", nil)
		resp, _ := app.Test(req)
		if reqID := resp.Header.Get("X-Request-ID"); reqID == "" {
			t.Errorf("X-Request-ID header tidak ditemukan pada response")
		}
	})
}

func TestApp_RequireJSON(t *testing.T) {
	app := config.NewApp()
	app.Post("/json-endpoint", middleware.RequireJSON, func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	t.Run("RejectNonJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json-endpoint", strings.NewReader("text payload"))
		req.Header.Set("Content-Type", "text/plain")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("AllowApplicationJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/json-endpoint", strings.NewReader(`{"key":"value"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})
}
