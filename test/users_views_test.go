package test

import (
	"bytes"
	"encoding/json"
	"go-crud/schemas"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	requestBody := map[string]string{
		"name":   "Connor Tran",
		"email": "connortran@gmail.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response schemas.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Connor Tran", response.Data.Name)
	assert.Equal(t, "connortran@gmail.com", response.Data.Email)
}

func TestCreateUserValidationError(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	requestBody := map[string]string{
		"name":   "",
		"email": "invalid-email",
		"password": "",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Validation failed")
}

func TestCreateUserMissingFields(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	requestBody := map[string]string{
		"name":   "Connor Tran",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Validation failed")
}

func TestGetUserByIDSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	user := UserFactory("testPassword123")
	req, _ := http.NewRequest("GET", "/users/"+strconv.FormatUint(uint64(user.ID), 10), nil)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, response.Data.Name)
	assert.Equal(t, user.Email, response.Data.Email)
}

func TestGetUserByIDNotFound(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	req, _ := http.NewRequest("GET", "/users/9999", nil)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "user not found")
}

func TestPartialUpdateUserSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	user := UserFactory("testPassword123",
		WithEmail("test-update@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"name": "Updated Name",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PATCH", "/users/"+strconv.FormatUint(uint64(user.ID), 10), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-update@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", response.Data.Name)
	assert.Equal(t, user.Email, response.Data.Email)
}

func TestPartialUpdateUserNotFound(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	UserFactory("testPassword123",
		WithEmail("test-auth@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"name": "Updated Name",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PATCH", "/users/9999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-auth@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "User not found")
}

func TestPartialUpdateUserValidationError(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	user := UserFactory("testPassword123",
		WithEmail("test-validation@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"email": "invalid-email",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PATCH", "/users/"+strconv.FormatUint(uint64(user.ID), 10), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-validation@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Validation failed")
}

func TestDeleteUserSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	user := UserFactory("testPassword123",
		WithEmail("test-delete@example.com"),
		WithName("Test User"),
	)

	req, _ := http.NewRequest("DELETE", "/users/"+strconv.FormatUint(uint64(user.ID), 10), nil)
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-delete@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response schemas.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "User deleted successfully")
}

func TestDeleteUserNotFound(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	UserFactory("testPassword123",
		WithEmail("test-delete-notfound@example.com"),
		WithName("Test User"),
	)

	req, _ := http.NewRequest("DELETE", "/users/9999", nil)
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-delete-notfound@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response schemas.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "User not found")
}
