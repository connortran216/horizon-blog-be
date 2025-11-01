package schemas

import (
	"go-crud/models"
)

// Query Parameters
type ListPostsQueryParams struct {
	Page     int      `json:"page" form:"page" validate:"omitempty,min=1" default:"1"`
	Limit    int      `json:"limit" form:"limit" validate:"omitempty,min=1,max=100" default:"10"`
	UserID   *uint    `json:"user_id,omitempty" form:"user_id"`
	TagNames []string `json:"tag_names,omitempty" form:"tags"` // Filter by tag names
}

// Input Schemas
type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required,min=1,max=255" example:"My New Post"`
	TagNames []string `json:"tag_names,omitempty" example:"golang,web-development,tutorial"`
}

// Method for CreatePostRequest struct
func (r CreatePostRequest) Validate() error {
	return validate.Struct(r)
}

type CreateVersionRequest struct {
	Title           string `json:"title" binding:"required,min=1,max=255" example:"My Post Version"`
	ContentMarkdown string `json:"content_markdown" binding:"required,min=1" example:"# My Post\n\n**markdown**"`
	ContentJSON     string `json:"content_json" binding:"required,min=1" example:"{\"type\":\"doc\",\"content\":[]}"`
}

// Method for CreateVersionRequest struct
func (r CreateVersionRequest) Validate() error {
	return validate.Struct(r)
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
