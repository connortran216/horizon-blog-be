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

func TestCreatePostSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	_ = UserFactory("testPassword123",
		WithEmail("test-update@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"title":            "Test Post Title",
		"content_markdown": "This is a test post content",
		"content_json":     "{\"type\":\"doc\",\"content\":[]}",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-update@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response schemas.PostResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Test Post Title", response.Data.Title)
	assert.Equal(t, "This is a test post content", response.Data.ContentMarkdown)
	assert.Equal(t, "Post created successfully", response.Message)
	assert.NotZero(t, response.Data.ID)
}

func TestCreatePostValidationError(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create user for authentication
	user := UserFactory("testPassword123",
		WithEmail("test-validation@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"title":            "", // Invalid - empty title
		"content_markdown": "This is a test post content",
		"content_json":     "{}",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, user.Email))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response.Error, "Invalid request data: Key: 'CreatePostRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag")
}

func TestGetPostByIDSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create mock data Post
	post := PostFactory()
	req, _ := http.NewRequest("GET", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.PostResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, post.Title, response.Data.Title)
	assert.Equal(t, post.ContentMarkdown, response.Data.ContentMarkdown)
}

func TestGetPostByIDDataDoesNotExist(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	req, _ := http.NewRequest("GET", "/posts/9999", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, response.Error, "post not found")
}

func TestListPostsSuccessWithDefaultPagination(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create mock data Posts
	for i := 0; i < 10; i++ {
		PostFactory()
	}

	req, _ := http.NewRequest("GET", "/posts", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ListPostsResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.GreaterOrEqual(t, len(response.Data), 10)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, 1, response.Page)
	assert.GreaterOrEqual(t, response.Total, 10)
}

func TestListPostsSuccessWithCustomPagination(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create mock data Posts
	for i := 0; i < 10; i++ {
		PostFactory()
	}

	req, _ := http.NewRequest("GET", "/posts?page=2&limit=5", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ListPostsResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.GreaterOrEqual(t, len(response.Data), 5)
	assert.Equal(t, 5, response.Limit)
	assert.Equal(t, 2, response.Page)
	assert.GreaterOrEqual(t, response.Total, 10)
}

func TestListPostsShouldReturnDefaultPaginationWhenInvalidQueryParams(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	req, _ := http.NewRequest("GET", "/posts?page=abc&limit=xyz", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ListPostsResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.GreaterOrEqual(t, len(response.Data), 0)
	assert.Equal(t, 10, response.Limit) // Default limit
	assert.Equal(t, 1, response.Page)   // Default page
	assert.GreaterOrEqual(t, response.Total, 0)
}

func TestUpdatePostSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create user and login to get auth token
	user := UserFactory("testPassword123",
		WithEmail("test-update@example.com"),
		WithName("Test User"),
	)

	// Create post for this user
	post := PostFactory(
		WithUserID(user.ID),
		WithTitle("Original Title"),
		WithContent("Original Content"),
	)

	requestBody := map[string]string{
		"title":            "Updated Title",
		"content_markdown": "Updated Content",
		"content_json":     "{\"type\":\"doc\",\"content\":[]}",
		"status":           "published",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-update@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.PostResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Updated Title", response.Data.Title)
	assert.Equal(t, "Updated Content", response.Data.ContentMarkdown)
	assert.NotZero(t, response.Data.ID)
}

func TestUpdatePostFailWhenDataDoesNotExist(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create authenticated user but try to update non-existent post
	UserFactory("testPassword123",
		WithEmail("test-update-nonexist@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"title":   "Updated Title",
		"content": "Updated Content",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/posts/9999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-update-nonexist@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, response.Error, "Post not found") // Match actual error message case
}

func TestUpdatePostFailWhenDataIsInvalid(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create authenticated user and post
	user := UserFactory("testPassword123",
		WithEmail("test-update-invalid@example.com"),
		WithName("Test User"),
	)
	post := PostFactory(WithUserID(user.ID))

	requestBody := map[string]string{
		"author":  "", // Invalid non-existent field
		"content": "Updated Content",
		"status":  "published",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-update-invalid@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, response.Error, "Invalid request data: Key: 'UpdatePostRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag")
}

func TestPartiallyUpdatePostSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create authenticated user and post for PATCH test
	user := UserFactory("testPassword123",
		WithEmail("test-patch@example.com"),
		WithName("Test User"),
	)
	post := PostFactory(
		WithUserID(user.ID),
		WithTitle("Original Title"),
		WithContent("Original Content"),
	)

	requestBody := map[string]string{
		"content_markdown": "Partially Updated Content",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PATCH", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-patch@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.PostResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Original Title", response.Data.Title) // Title should remain unchanged
	assert.Equal(t, "Partially Updated Content", response.Data.ContentMarkdown)
	assert.NotZero(t, response.Data.ID)
}

func TestPartiallyUpdatePostFailWhenDataDoesNotExist(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create authenticated user but try to patch non-existent post
	UserFactory("testPassword123",
		WithEmail("test-patch-nonexist@example.com"),
		WithName("Test User"),
	)

	requestBody := map[string]string{
		"content_markdown": "Partially Updated Content",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PATCH", "/posts/9999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-patch-nonexist@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, response.Error, "Post not found") // Match actual error message case
}

func TestPartiallyUpdatePostFailInvalidData(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create authenticated user and post
	user := UserFactory("testPassword123",
		WithEmail("test-patch-invalid@example.com"),
		WithName("Test User"),
	)
	post := PostFactory(WithUserID(user.ID))

	requestBody := map[string]string{
		"title": "", // Invalid empty title
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PATCH", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-patch-invalid@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, response.Error, "Invalid request data: Key: 'PatchPostRequest.Title' Error:Field validation for 'Title' failed on the 'min' tag")
}

func TestDeletePostSuccess(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create authenticated user and post for DELETE test
	user := UserFactory("testPassword123",
		WithEmail("test-delete@example.com"),
		WithName("Test User"),
	)
	post := PostFactory(
		WithUserID(user.ID),
		WithTitle("Post to delete"),
		WithContent("This post will be deleted"),
	)

	req, _ := http.NewRequest("DELETE", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-delete@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.MessageResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Post deleted successfully", response.Message)
}

func TestDeletePostFailWhenDataDoesNotExist(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create authenticated user but try to delete non-existent post
	UserFactory("testPassword123",
		WithEmail("test-delete-nonexist@example.com"),
		WithName("Test User"),
	)

	req, _ := http.NewRequest("DELETE", "/posts/9999", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "test-delete-nonexist@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, response.Error, "Post not found") // Match actual error message case
}

// Additional tests for ownership restrictions
func TestUpdatePostFailWhenWrongUser(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create first user and their post
	user1 := UserFactory("testPassword123",
		WithEmail("user1@example.com"),
		WithName("User One"),
	)
	user1Post := PostFactory(WithUserID(user1.ID))

	// Create second user who tries to update the post
	UserFactory("testPassword123",
		WithEmail("user2@example.com"),
		WithName("User Two"),
	)

	requestBody := map[string]string{
		"title":   "Hacked Title",
		"content": "Hacked Content",
		"status":  "published",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/posts/"+strconv.FormatUint(uint64(user1Post.ID), 10), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "user2@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, response.Error, "You can only update your own posts")
}

func TestDeletePostFailWhenWrongUser(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.TearDown()

	// Create first user and their post
	user1 := UserFactory("testPassword123",
		WithEmail("user1-delete@example.com"),
		WithName("User One"),
	)
	post := PostFactory(WithUserID(user1.ID))

	// Create second user who tries to delete the post
	_ = UserFactory("testPassword123",
		WithEmail("user2-delete@example.com"),
		WithName("User Two"),
	)

	req, _ := http.NewRequest("DELETE", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "user2-delete@example.com"))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response schemas.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, response.Error, "You can only delete your own posts")
}

// func TestListPostsWithMineFilter(t *testing.T) {
// 	suite := NewTestSuite(t)
// 	defer suite.TearDown()

// 	// Create first user with 2 posts
// 	user1 := UserFactory("testPassword123",
// 		WithEmail("user1-mine@example.com"),
// 		WithName("User One"),
// 	)
// 	_ = PostFactory(WithUserID(user1.ID), WithTitle("User1 Post 1"))
// 	_ = PostFactory(WithUserID(user1.ID), WithTitle("User1 Post 2"))

// 	// Create second user with 1 post
// 	user2 := UserFactory("testPassword123",
// 		WithEmail("user2-mine@example.com"),
// 		WithName("User Two"),
// 	)
// 	_ = PostFactory(WithUserID(user2.ID), WithTitle("User2 Post 1"))

// 	// Login as user1 and filter by "mine"
// 	req, _ := http.NewRequest("GET", "/posts?mine=true", nil)
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+getAuthToken(t, suite, "user1-mine@example.com"))

// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	var response schemas.ListPostsResponse
// 	json.Unmarshal(w.Body.Bytes(), &response)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Equal(t, 2, len(response.Data)) // Should only see user1's posts

// 	// Verify post titles
// 	postTitles := make([]string, len(response.Data))
// 	for i, post := range response.Data {
// 		postTitles[i] = post.Title
// 	}
// 	assert.Contains(t, postTitles, "User1 Post 1")
// 	assert.Contains(t, postTitles, "User1 Post 2")
// 	assert.NotContains(t, postTitles, "User2 Post 1")
// }
