package helper_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jamaalul/proyek1pbl/helper"
)

func TestAppError_Constructors(t *testing.T) {
	tests := []struct {
		name       string
		err        *helper.AppError
		wantStatus int
		wantMsg    string
	}{
		{"BadRequest", helper.BadRequest("bad request"), http.StatusBadRequest, "bad request"},
		{"Unauthorized", helper.Unauthorized("unauthorized"), http.StatusUnauthorized, "unauthorized"},
		{"Forbidden", helper.Forbidden("forbidden"), http.StatusForbidden, "forbidden"},
		{"NotFound", helper.NotFound("not found"), http.StatusNotFound, "not found"},
		{"Conflict", helper.Conflict("conflict"), http.StatusConflict, "conflict"},
		{"TooManyRequests", helper.TooManyRequests("too many"), http.StatusTooManyRequests, "too many"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Status != tt.wantStatus {
				t.Errorf("%s status = %d, want %d", tt.name, tt.err.Status, tt.wantStatus)
			}
			if tt.err.Message != tt.wantMsg {
				t.Errorf("%s message = %s, want %s", tt.name, tt.err.Message, tt.wantMsg)
			}
		})
	}

	t.Run("Validation", func(t *testing.T) {
		errs := map[string][]string{"email": {"email tidak valid"}}
		valErr := helper.Validation(errs)
		if valErr.Status != http.StatusUnprocessableEntity {
			t.Errorf("Validation status = %d, want %d", valErr.Status, http.StatusUnprocessableEntity)
		}
		if valErr.Errors["email"][0] != "email tidak valid" {
			t.Errorf("Validation errors map tidak sesuai: %+v", valErr.Errors)
		}
	})

	t.Run("InternalAndUnwrap", func(t *testing.T) {
		cause := errors.New("koneksi database terputus")
		intErr := helper.Internal(cause)
		if intErr.Status != http.StatusInternalServerError {
			t.Errorf("Internal status = %d, want %d", intErr.Status, http.StatusInternalServerError)
		}
		if !errors.Is(intErr, cause) {
			t.Errorf("errors.Is gagal melakukan unwrap pada cause: %v", intErr)
		}
		if intErr.Error() != "terjadi kesalahan pada server: koneksi database terputus" {
			t.Errorf("Error() string format tidak sesuai: %s", intErr.Error())
		}
	})
}
