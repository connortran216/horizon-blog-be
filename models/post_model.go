package models

import "time"

type PostStatus string

const (
    Draft   PostStatus = "draft"
    Published   PostStatus = "published"
)


type Post struct {
	ID              uint       `gorm:"primaryKey" json:"id" example:"1"`
	UserID          uint       `gorm:"not null" json:"user_id" example:"1"`
	Title           string     `gorm:"not null" json:"title" example:"My First Post"`
	ContentMarkdown string     `gorm:"column:content_markdown;type:text" json:"content_markdown" example:"# My First Post\n\nThis is **markdown** content"`
	ContentJSON     string     `gorm:"column:content_json;type:text" json:"content_json" example:"{\"type\":\"doc\",\"content\":[]}"`
	Status          PostStatus `gorm:"default:'draft';not null" json:"status" example:"draft"`
	User            *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tags            []Tag      `gorm:"many2many:post_tags" json:"tags,omitempty"`
	CreatedAt       time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt       time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// GetID implements the ModelInterface
func (p Post) GetID() uint {
	return p.ID
}
