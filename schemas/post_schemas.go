package schemas

import (
	"go-crud/models"
	"github.com/go-playground/validator/v10"
)


var validate = validator.New()

// Query Parameters
type ListPostsQueryParams struct {
	Page   int  `json:"page" form:"page" validate:"omitempty,min=1" default:"1"`
	Limit  int  `json:"limit" form:"limit" validate:"omitempty,min=1,max=100" default:"10"`
	UserID *uint `json:"user_id,omitempty" form:"user_id"`
}

// Input Schemas
type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255" example:"My New Post"`
	Content string `json:"content" validate:"required,min=1" example:"This is the content of my new post"`
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
	Title   string `json:"title" validate:"required,min=1,max=255" example:"Updated Post Title"`
	Content string `json:"content" validate:"required,min=1" example:"Updated post content"`
}

// Method for UpdatePostRequest struct
func (r UpdatePostRequest) Validate() error {
	return validate.Struct(r)
}

func (r UpdatePostRequest) ToModel() models.Post {
	return models.Post{
		Title:   r.Title,
		Content: r.Content,
	}
}

type PatchPostRequest struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,min=1,max=255" example:"Partially Updated Title"`
	Content *string `json:"content,omitempty" validate:"omitempty,min=1" example:"Partially updated content"`
}

// Method for PatchPostRequest struct
func (r PatchPostRequest) Validate() error {
	return validate.Struct(r)
}

func (r PatchPostRequest) IsEmpty() bool {
	return r.Title == nil && r.Content == nil
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
