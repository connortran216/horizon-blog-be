package schemas

import (
	"go-crud/models"
)

// Query Parameters
type ListPostsQueryParams struct {
	Page     int    `json:"page" form:"page" validate:"omitempty,min=1" default:"1"`
	Limit    int    `json:"limit" form:"limit" validate:"omitempty,min=1,max=100" default:"10"`
	UserID   *uint  `json:"user_id,omitempty" form:"user_id"`
	TagNames []string `json:"tag_names,omitempty" form:"tags"` // Filter by tag names
}

// Input Schemas
type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required,min=1,max=255" example:"My New Post"`
	Content  string   `json:"content" binding:"required,min=1" example:"This is the content of my new post"`
	TagNames []string `json:"tag_names,omitempty" example:"golang,web-development,tutorial"`
}

// Method for CreatePostRequest struct
func (r CreatePostRequest) Validate() error {
	return validate.Struct(r)
}

func (r CreatePostRequest) ToModel() models.Post {
	return models.Post{
		Title:   r.Title,
		Content: r.Content,
	}
}

type UpdatePostRequest struct {
	Title   string           `json:"title" binding:"required,min=1,max=255" example:"Updated Post Title"`
	Content string           `json:"content" binding:"required,min=1" example:"Updated post content"`
	Status  models.PostStatus `json:"status" binding:"required,oneof=draft published" example:"published"`
}

// Method for UpdatePostRequest struct
func (r UpdatePostRequest) Validate() error {
	return validate.Struct(r)
}

func (r UpdatePostRequest) ToModel() models.Post {
	return models.Post{
		Title:   r.Title,
		Content: r.Content,
		Status:  r.Status,
	}
}

type PatchPostRequest struct {
	Title   *string              `json:"title,omitempty" binding:"omitempty,min=1,max=255" example:"Partially Updated Title"`
	Content *string              `json:"content,omitempty" binding:"omitempty,min=1" example:"Partially updated content"`
	Status  *models.PostStatus `json:"status,omitempty" binding:"omitempty,oneof=draft published" example:"published"`
}

// Method for PatchPostRequest struct
func (r PatchPostRequest) Validate() error {
	return validate.Struct(r)
}

func (r PatchPostRequest) IsEmpty() bool {
	return r.Title == nil && r.Content == nil && r.Status == nil
}

// Method for PatchPostRequest struct
func (r PatchPostRequest) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	if r.Title != nil {
		data["title"] = *r.Title
	}
	if r.Content != nil {
		data["content"] = *r.Content
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
