package views

import (
	"fmt"
	"go-crud/schemas"
	"go-crud/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagViews struct {
	service *services.TagService
}

func NewTagViews() *TagViews {
	return &TagViews{
		service: services.NewTagService(),
	}
}

// @Summary Create tag
// @Tags tags
// @Param tag body schemas.CreateTagRequest true "Tag data"
// @Success 201 {object} schemas.TagResponse
// @Router /tags [post]
func (v *TagViews) CreateTag(c *gin.Context) {
	var input schemas.CreateTagRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Invalid request data: %v", err),
		})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Validation failed: %v", err),
		})
		return
	}

	tagModel := input.ToModel()
	result, err := v.service.Create(tagModel)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "tag already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to create tag: %v", err),
		})
		return
	}

	response := schemas.TagResponse{
		Data:    *result,
		Message: "Tag created successfully",
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary List tags
// @Tags tags
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort query string false "Sort by" Enums(name, usage_count, created_at) default(name)
// @Success 200 {object} schemas.ListTagsResponse
// @Router /tags [get]
func (v *TagViews) ListTags(c *gin.Context) {
	var query schemas.TagListQueryParams
	query.Page, _ = strconv.Atoi(c.Query("page"))
	query.Limit, _ = strconv.Atoi(c.Query("limit"))
	query.Sort = c.Query("sort")

	// Set defaults if not provided
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	}
	if query.Sort == "" {
		query.Sort = "name"
	}

	results, total, err := v.service.GetAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to fetch tags: %v", err),
		})
		return
	}

	response := schemas.ListTagsResponse{
		Data:  results,
		Total: int(total),
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Get tag
// @Tags tags
// @Param id path int true "Tag ID"
// @Success 200 {object} schemas.TagResponse
// @Router /tags/{id} [get]
func (v *TagViews) GetTag(c *gin.Context) {
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
			Error: fmt.Sprintf("Tag not found: %v", err),
		})
		return
	}

	response := schemas.TagResponse{
		Data: *result,
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Get popular tags
// @Tags tags
// @Param limit query int false "Limit number of tags" default(10)
// @Success 200 {object} schemas.ListTagsResponse
// @Router /tags/popular [get]
func (v *TagViews) GetPopularTags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}

	results, err := v.service.GetPopular(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to fetch popular tags: %v", err),
		})
		return
	}

	response := schemas.ListTagsResponse{
		Data:  results,
		Total: len(results),
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Search tags
// @Tags tags
// @Param q query string true "Search query"
// @Param limit query int false "Limit number of results" default(10)
// @Success 200 {object} schemas.ListTagsResponse
// @Router /tags/search [get]
func (v *TagViews) SearchTags(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Search query is required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}

	results, err := v.service.Search(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to search tags: %v", err),
		})
		return
	}

	response := schemas.ListTagsResponse{
		Data:  results,
		Total: len(results),
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Update tag
// @Tags tags
// @Param id path int true "Tag ID"
// @Param tag body schemas.UpdateTagRequest true "Tag data"
// @Success 200 {object} schemas.TagResponse
// @Router /tags/{id} [put]
func (v *TagViews) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	var input schemas.UpdateTagRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Invalid request data: %v", err),
		})
		return
	}

	result, err := v.service.Update(uint(id), input.ToMap())
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "tag not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to update tag: %v", err),
		})
		return
	}

	response := schemas.TagResponse{
		Data:    *result,
		Message: "Tag updated successfully",
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Delete tag
// @Tags tags
// @Param id path int true "Tag ID"
// @Success 200 {object} schemas.MessageResponse
// @Router /tags/{id} [delete]
func (v *TagViews) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	err = v.service.Delete(uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "tag not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to delete tag: %v", err),
		})
		return
	}

	response := schemas.MessageResponse{
		Message: "Tag deleted successfully",
	}
	c.JSON(http.StatusOK, response)
}

func (v *TagViews) RegisterRoutes(router *gin.Engine) {
	tags := router.Group("/tags")
	{
		tags.POST("", AuthMiddleware(), v.CreateTag)
		tags.GET("", v.ListTags)
		tags.GET("/popular", v.GetPopularTags)
		tags.GET("/search", v.SearchTags)
		tags.GET("/:id", v.GetTag)
		tags.PUT("/:id", AuthMiddleware(), v.UpdateTag)
		tags.DELETE("/:id", AuthMiddleware(), v.DeleteTag)
	}
}
