package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

// TestCreateParticipant_Validation verifies the HTTP boundary validations.
func TestCreateParticipant_Validation(t *testing.T) {
	app := fiber.New()
	
	// Create handler with nil services, as we are only testing
	// the early boundary checks (validation) which occur before
	// any service calls.
	h := NewHandler(nil, nil, nil, nil)
	app.Post("/participants", h.CreateParticipant)

	tests := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "Empty Payload",
			payload:      nil,
			expectedCode: 400,
		},
		{
			name:         "Missing Name",
			payload:      models.ParticipantRequest{Age: 20, Gender: "male"},
			expectedCode: 400,
		},
		{
			name:         "Age Under 15",
			payload:      models.ParticipantRequest{Name: "Test", Age: 12, Gender: "male"},
			expectedCode: 400,
		},
		{
			name:         "Invalid Gender",
			payload:      models.ParticipantRequest{Name: "Test", Age: 25, Gender: "alien"},
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/participants", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}
