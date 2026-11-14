package model

// WebResponse adalah format standar JSON response API
type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}
