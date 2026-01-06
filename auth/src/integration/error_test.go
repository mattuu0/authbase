package integration

import (
	"auth/models"
	"auth/services"
	testtool "auth/testing"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupTestRouter()

	t.Run("invalid credentials", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Attempting login with invalid credentials")
		loginBody := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "WrongPassword123!",
		}
		loginJSON, _ := json.Marshal(loginBody)

		req := httptest.NewRequest(http.MethodPost, "/basic/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Login failed as expected")
	})

	t.Run("malformed request body", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Sending malformed JSON")
		req := httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		testtool.LogSuccess(t, "Request rejected as bad request")
	})

	t.Run("missing authorization header", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Accessing protected endpoint without token")
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		testtool.LogSuccess(t, "Request rejected as unauthorized")
	})
}
