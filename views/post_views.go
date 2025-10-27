package views

import (
	"fmt"
	"go-crud/initializers"
	"go-crud/models"
	"go-crud/schemas"
	"go-crud/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PostViews struct {
	service *services.PostService
}

func NewPostViews() *PostViews {
	return &PostViews{
		service: services.NewPostService(),
	}
}

// @Summary Create post
// @Tags posts
// @Param post body schemas.CreatePostRequest true "Post data"
// @Success 201 {object} schemas.PostResponse
// @Router /posts [post]
func (v *PostViews) CreatePost(c *gin.Context) {
	var input schemas.CreatePostRequest
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

	// Create post metadata
	result, err := v.service.Create(models.Post{
		UserID: userID,
		Title:  input.Title,
	}, input.TagNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to create post: %v", err),
		})
		return
	}

	response := schemas.PostResponse{
		Data:    *result,
		Message: "Post created successfully",
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary List posts
// @Tags posts
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} schemas.ListPostsResponse
// @Router /posts [get]
func (v *PostViews) ListPosts(c *gin.Context) {
	var query schemas.ListPostsQueryParams
	query.Page, _ = strconv.Atoi(c.Query("page"))
	query.Limit, _ = strconv.Atoi(c.Query("limit"))

	// Set defaults if not provided
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	// Check if authenticated user wants to filter by their own posts
	userID, exists := GetUserIDFromContext(c)
	if exists {
		if c.Query("mine") == "true" {
			query.UserID = &userID
		}
	}

	// Handle tag filtering
	if tagsParam := c.Query("tags"); tagsParam != "" {
		// Split comma-separated tag names
		tagNames := strings.Split(tagsParam, ",")
		// Trim spaces and filter out empty strings
		for i, tag := range tagNames {
			tagNames[i] = strings.TrimSpace(tag)
		}
		// Filter out empty tag names
		var filteredTags []string
		for _, tag := range tagNames {
			if tag != "" {
				filteredTags = append(filteredTags, tag)
			}
		}
		query.TagNames = filteredTags
	}

	results, total, err := v.service.GetWithPagination(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to fetch posts: %v", err),
		})
		return
	}

response := schemas.ListPostsResponse{
		Data:  results,
		Limit: query.Limit,
		Page:  query.Page,
		Total: int(total),
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Get post
// @Tags posts
// @Param id path int true "Post ID"
// @Success 200 {object} schemas.PostResponse
// @Router /posts/{id} [get]
func (v *PostViews) GetPost(c *gin.Context) {
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
			Error: fmt.Sprintf("Post not found: %v", err),
		})
		return
	}

	response := schemas.PostResponse{
		Data: *result,
	}
	c.JSON(http.StatusOK, response)
}



// @Summary Delete post
// @Tags posts
// @Param id path int true "Post ID"
// @Success 200 {object} schemas.MessageResponse
// @Router /posts/{id} [delete]
func (v *PostViews) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	// Check if post exists and belongs to authenticated user
	authenticatedUserID, exists := GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	post, err := v.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: "Post not found",
		})
		return
	}

	if post.UserID != authenticatedUserID {
		c.JSON(http.StatusForbidden, schemas.ErrorResponse{
			Error: "You can only delete your own posts",
		})
		return
	}

	err = v.service.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to delete post: %v", err),
		})
		return
	}

	response := schemas.MessageResponse{
		Message: "Post deleted successfully",
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Update post
// @Tags posts
// @Param id path int true "Post ID"
// @Param post body schemas.UpdatePostRequest true "Post data"
// @Success 200 {object} schemas.PostResponse
// @Router /posts/{id} [put]
func (v *PostViews) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	// Check if post exists and belongs to authenticated user
	authenticatedUserID, exists := GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	post, err := v.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: "Post not found",
		})
		return
	}

	if post.UserID != authenticatedUserID {
		c.JSON(http.StatusForbidden, schemas.ErrorResponse{
			Error: "You can only update your own posts",
		})
		return
	}

	var input schemas.CreatePostRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Invalid request data: %v", err),
		})
		return
	}

	// Basic validation for title
	if input.Title == "" {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid request data: Key: 'CreatePostRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag",
		})
		return
	}

	// In new architecture, only allow title updates for posts
	// Content updates are handled via versions
	post.Title = input.Title
	result := initializers.DB.Save(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: "Failed to update post",
		})
		return
	}

	response := schemas.PostResponse{
		Data:    *post,
		Message: "Post updated successfully",
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Partially update post
// @Tags posts
// @Param id path int true "Post ID"
// @Param post body schemas.PatchPostRequest true "Partial post data"
// @Success 200 {object} schemas.PostResponse
// @Router /posts/{id} [patch]
func (v *PostViews) PartialUpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	// Check if post exists and belongs to authenticated user
	authenticatedUserID, exists := GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	post, err := v.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: "Post not found",
		})
		return
	}

	if post.UserID != authenticatedUserID {
		c.JSON(http.StatusForbidden, schemas.ErrorResponse{
			Error: "You can only update your own posts",
		})
		return
	}

	var input schemas.PatchPostRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Apply partial updates
	if input.Title != nil {
		if *input.Title == "" {
			c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
				Error: "Title cannot be empty",
			})
			return
		}
		post.Title = *input.Title
	}

	// Note: In new architecture, status and content are managed via versions
	result := initializers.DB.Save(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: "Failed to update post",
		})
		return
	}

	response := schemas.PostResponse{
		Data:    *post,
		Message: "Post updated successfully",
	}
	c.JSON(http.StatusOK, response)
}

func (v *PostViews) RegisterRoutes(router *gin.Engine) {
	posts := router.Group("/posts")
	{
		posts.POST("", AuthMiddleware(), v.CreatePost)
		posts.GET("", v.ListPosts)
		posts.GET("/:id", v.GetPost)
		posts.PUT("/:id", AuthMiddleware(), v.UpdatePost)
		posts.PATCH("/:id", AuthMiddleware(), v.PartialUpdatePost)
		posts.DELETE("/:id", AuthMiddleware(), v.DeletePost)
	}
}
