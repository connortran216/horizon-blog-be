package views

import (
	"fmt"
	"go-crud/models"
	"go-crud/schemas"
	"go-crud/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostVersionViews struct {
	service *services.PostVersionService
}

func NewPostVersionViews() *PostVersionViews {
	return &PostVersionViews{
		service: services.NewPostVersionService(),
	}
}

// @Summary List post versions globally
// @Tags versions
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param user_id query int false "Filter by user ID"
// @Param status query string false "Filter by status (draft, published)"
// @Success 200 {object} schemas.ListPostsResponse // Reuse for versions, or create ListVersionsResponse
// @Router /versions [get]
func (v *PostVersionViews) ListVersions(c *gin.Context) {
	// TODO: Implement with pagination and filtering
	c.JSON(http.StatusNotImplemented, gin.H{"message": "ListVersions not implemented"})
}

// @Summary Get version by ID
// @Tags versions
// @Param id path int true "Version ID"
// @Success 200 {object} schemas.PostResponse // Reuse or create VersionResponse
// @Router /versions/{id} [get]
func (v *PostVersionViews) GetVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	result, err := v.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: fmt.Sprintf("Version not found: %v", err),
		})
		return
	}

	// TODO: Create or use appropriate response schema
	response := schemas.PostResponse{
		Data: *result.Post, // Hack, reuse
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Create draft version for a post
// @Tags versions
// @Param id path int true "Post ID"
// @Param version body schemas.CreateVersionRequest true "Version content"
// @Success 201 {object} schemas.PostResponse
// @Router /posts/{id}/versions [post]
func (v *PostVersionViews) CreateVersion(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid post ID format",
		})
		return
	}

	var input schemas.CreateVersionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Invalid request data: %v", err),
		})
		return
	}

	// Get authenticated user ID
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	result, err := v.service.CreateDraftVersion(uint(postID), userID, input.Title, input.ContentMarkdown, input.ContentJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to create version: %v", err),
		})
		return
	}

	response := schemas.PostResponse{
		Data: *result.Post, // Hack
		Message: "Version created successfully",
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary List versions for a specific post
// @Tags versions
// @Param id path int true "Post ID"
// @Success 200 {object} schemas.ListPostsResponse
// @Router /posts/{id}/versions [get]
func (v *PostVersionViews) ListVersionsForPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid post ID format",
		})
		return
	}

	results, err := v.service.GetVersionsForPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to fetch versions: %v", err),
		})
		return
	}

	// Transform to posts for reuse
	posts := make([]models.Post, len(results))
	for i, version := range results {
		post := *version.Post
		post.Title = version.Title // Hack
		posts[i] = post
	}

	response := schemas.ListPostsResponse{
		Data:  posts, // TODO: ListVersionsResponse
		Limit: 100,   // Not paginated
		Page:  1,
		Total: len(results),
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Publish a version
// @Tags versions
// @Param id path int true "Version ID"
// @Success 200 {object} schemas.PostResponse
// @Router /versions/{id}/publish [patch]
func (v *PostVersionViews) PublishVersion(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid version ID format",
		})
		return
	}

	userID, exists := GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	result, err := v.service.PublishVersion(uint(versionID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to publish version: %v", err),
		})
		return
	}

	response := schemas.PostResponse{
		Data: *result.Post,
		Message: "Version published successfully",
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Update a version (auto-save draft)
// @Tags versions
// @Param id path int true "Version ID"
// @Param version body schemas.UpdatePostRequest true "Version content"
// @Success 200 {object} schemas.PostResponse
// @Router /versions/{id} [put]
func (v *PostVersionViews) UpdateVersion(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid version ID format",
		})
		return
	}

	var input schemas.CreateVersionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Invalid request data: %v", err),
		})
		return
	}

	userID, exists := GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	result, err := v.service.AutoSaveDraft(uint(versionID), userID, input.Title, input.ContentMarkdown, input.ContentJSON)
	if err != nil {
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to update version: %v", err),
		})
		return
	}

	response := schemas.PostResponse{
		Data: *result.Post,
		Message: "Version updated successfully",
	}
	c.JSON(http.StatusOK, response)
}

func (v *PostVersionViews) RegisterRoutes(router *gin.Engine) {
	versions := router.Group("/versions")
	{
		versions.GET("", v.ListVersions)
		versions.GET("/:id", v.GetVersion)
		versions.PUT("/:id", AuthMiddleware(), v.UpdateVersion)
		versions.PATCH("/:id/publish", AuthMiddleware(), v.PublishVersion)
	}

	// Nested under posts
	posts := router.Group("/posts")
	{
		posts.GET("/:id/versions", v.ListVersionsForPost)
		posts.POST("/:id/versions", AuthMiddleware(), v.CreateVersion)
	}
}
