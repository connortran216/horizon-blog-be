package schemas

import (
	"go-crud/models"
)

// Tag Schemas
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50" example:"golang"`
	Description string `json:"description" binding:"max=500" example:"Go programming language"`
}

func (r CreateTagRequest) Validate() error {
	return validate.Struct(r)
}

func (r CreateTagRequest) ToModel() models.Tag {
	return models.Tag{
		Name:        r.Name,
		Description: r.Description,
	}
}

type UpdateTagRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=50" example:"go-programming"`
	Description *string `json:"description" binding:"omitempty,max=500" example:"Updated description"`
}

func (r UpdateTagRequest) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	if r.Name != nil {
		data["name"] = *r.Name
	}
	if r.Description != nil {
		data["description"] = *r.Description
	}
	return data
}

// Response Schemas
type TagResponse struct {
	Data    models.Tag `json:"data"`
	Message string     `json:"message,omitempty"`
}

type ListTagsResponse struct {
	Data  []models.Tag `json:"data"`
	Total int          `json:"total"`
}

type TagListQueryParams struct {
	Page  int `json:"page" form:"page" validate:"omitempty,min=1" default:"1"`
	Limit int `json:"limit" form:"limit" validate:"omitempty,min=1,max=100" default:"20"`
	Sort  string `json:"sort" form:"sort" validate:"omitempty,oneof=name usage_count created_at" default:"name"`
}
