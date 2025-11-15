package schemas

import (
	"go-crud/models"
)

// Query Parameters
type ListPostsQueryParams struct {
	Page     int                `form:"page" binding:"omitempty,min=0"`
	Limit    int                `form:"limit" binding:"omitempty,min=0,max=100"`
	UserID   *uint              `form:"user_id"`
	Status   *models.PostStatus `form:"status"`  // Filter by status
	TagNames []string           `form:"tags"`    // Filter by tag names
}

// Method for ListPostsQueryParams struct - sets default values
func (q *ListPostsQueryParams) SetDefaults() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Status == nil {
		defaultStatus := models.Published
		q.Status = &defaultStatus
	}
}

// Input Schemas
type CreatePostRequest struct {
	Title           string             `json:"title" binding:"required,min=1,max=255" example:"My New Post"`
	ContentMarkdown string             `json:"content_markdown" binding:"required,min=1" example:"# My Post\n\nThis is **markdown**"`
	ContentJSON     string             `json:"content_json" binding:"required,min=1" example:"{\"type\":\"doc\",\"content\":[]}"`
	Status          *models.PostStatus `json:"status,omitempty" binding:"omitempty,oneof=draft published" example:"draft"`
	TagNames        []string           `json:"tag_names,omitempty" example:"golang,web-development,tutorial"`
}

// Method for CreatePostRequest struct
func (r CreatePostRequest) Validate() error {
	return validate.Struct(r)
}

func (r CreatePostRequest) ToModel() models.Post {
	post := models.Post{
		Title:           r.Title,
		ContentMarkdown: r.ContentMarkdown,
		ContentJSON:     r.ContentJSON,
	}

	// If status is provided, use it; otherwise default to draft
	if r.Status != nil {
		post.Status = *r.Status
	} else {
		post.Status = models.Draft
	}

	return post
}

type UpdatePostRequest struct {
	Title           string            `json:"title" binding:"required,min=1,max=255" example:"Updated Post Title"`
	ContentMarkdown string            `json:"content_markdown" binding:"required,min=1" example:"# Updated\n\nMarkdown content"`
	ContentJSON     string            `json:"content_json" binding:"required,min=1" example:"{\"type\":\"doc\",\"content\":[]}"`
	Status          models.PostStatus `json:"status" binding:"required,oneof=draft published" example:"published"`
}

// Method for UpdatePostRequest struct
func (r UpdatePostRequest) Validate() error {
	return validate.Struct(r)
}

func (r UpdatePostRequest) ToModel() models.Post {
	return models.Post{
		Title:           r.Title,
		ContentMarkdown: r.ContentMarkdown,
		ContentJSON:     r.ContentJSON,
		Status:          r.Status,
	}
}

type PatchPostRequest struct {
	Title           *string            `json:"title,omitempty" binding:"omitempty,min=1,max=255" example:"Partially Updated Title"`
	ContentMarkdown *string            `json:"content_markdown,omitempty" binding:"omitempty,min=1" example:"# Updated\n\nPartial markdown"`
	ContentJSON     *string            `json:"content_json,omitempty" binding:"omitempty,min=1" example:"{\"type\":\"doc\",\"content\":[]}"`
	Status          *models.PostStatus `json:"status,omitempty" binding:"omitempty,oneof=draft published" example:"published"`
}

// Method for PatchPostRequest struct
func (r PatchPostRequest) Validate() error {
	return validate.Struct(r)
}

func (r PatchPostRequest) IsEmpty() bool {
	return r.Title == nil && r.ContentMarkdown == nil && r.ContentJSON == nil && r.Status == nil
}

// Method for PatchPostRequest struct
func (r PatchPostRequest) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	if r.Title != nil {
		data["title"] = *r.Title
	}
	if r.ContentMarkdown != nil {
		data["content_markdown"] = *r.ContentMarkdown
	}
	if r.ContentJSON != nil {
		data["content_json"] = *r.ContentJSON
	}
	if r.Status != nil {
		data["status"] = *r.Status
	}
	return data
}

// Output Schemas
type PostResponse struct {
	Data    models.Post `json:"data"`
	Message string      `json:"message,omitempty"`
}

type ListPostsResponse struct {
	Data  []models.Post `json:"data"`
	Limit int           `json:"limit"`
	Page  int           `json:"page"`
	Total int           `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
