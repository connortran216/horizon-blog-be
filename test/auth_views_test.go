package test

import (
	"bytes"
	"encoding/json"
	"go-crud/schemas"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Set JWT secret for testing
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Setenv("JWT_SECRET", originalSecret)

	// Create a test user
	testUser := UserFactory("testPassword123",
		WithEmail("test@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"email":    "test@example.com",
		"password": "testPassword123", // Use the plain password
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, testUser.ID, response.User.Data.ID)
	assert.Equal(t, testUser.Email, response.User.Data.Email)
	assert.Equal(t, testUser.Name, response.User.Data.Name)
}

func TestLoginInvalidEmail(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	requestBody := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Invalid email or password")
}

func TestLoginInvalidPassword(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create a test user
	UserFactory("testPassword123",
		WithEmail("test@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"email":    "test@example.com",
		"password": "wrongPassword",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Invalid email or password")
}

func TestLoginValidationError(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	requestBody := map[string]string{
		"email":    "",
		"password": "",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Validation failed")
}

func TestLoginInvalidEmailFormat(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	requestBody := map[string]string{
		"email":    "invalid-email",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Validation failed")
}

// Helper function to get JWT token for authenticated requests
func getAuthToken(t *testing.T, suite *BaseTestSuite, email string) string {
	requestBody := map[string]string{
		"email":    email,
		"password": "testPassword123",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	return response.Token
}
