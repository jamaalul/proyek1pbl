package helper_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/app/model"
	"github.com/jamaalul/proyek1pbl/helper"
)

func TestResponse_Helpers(t *testing.T) {
	app := fiber.New()

	app.Get("/test-success", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "data berhasil diambil", fiber.Map{"name": "test"})
	})

	app.Get("/test-list", func(c *fiber.Ctx) error {
		meta := &model.Meta{CurrentPage: 1, PerPage: 10, Total: 50, LastPage: 5}
		return helper.SuccessList(c, "daftar data", []string{"a", "b"}, meta)
	})

	app.Post("/test-created", func(c *fiber.Ctx) error {
		return helper.Created(c, "data berhasil dibuat", fiber.Map{"id": 1}, "/items/1")
	})

	app.Delete("/test-nocontent", func(c *fiber.Ctx) error {
		return helper.NoContent(c)
	})

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-success", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request gagal: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var res model.WebResponse
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &res); err != nil {
			t.Fatalf("gagal unmarshal JSON response: %v", err)
		}
		if !res.Success || res.Message != "data berhasil diambil" {
			t.Errorf("response envelope tidak sesuai: %+v", res)
		}
	})

	t.Run("SuccessListWithMeta", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-list", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request gagal: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var res model.WebResponse
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &res); err != nil {
			t.Fatalf("gagal unmarshal: %v", err)
		}
		if res.Meta == nil {
			t.Errorf("meta tidak boleh nil pada SuccessList")
		}
	})

	t.Run("CreatedWithLocationHeader", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test-created", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request gagal: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
		}
		if loc := resp.Header.Get("Location"); loc != "/items/1" {
			t.Errorf("header Location = %s, want /items/1", loc)
		}
	})

	t.Run("NoContent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/test-nocontent", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("request gagal: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
		}
	})
}
