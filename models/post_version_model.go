package models

import "time"

type PostVersionStatus string

const (
	VersionDraft     PostVersionStatus = "draft"
	VersionPublished PostVersionStatus = "published"
)

type PostVersion struct {
	ID              uint              `gorm:"primaryKey" json:"id" example:"1"`
	PostID          uint              `gorm:"not null" json:"post_id" example:"1"`
	Title           string            `gorm:"not null" json:"title" example:"My Post Version"`
	ContentMarkdown string            `gorm:"column:content_markdown;type:text" json:"content_markdown" example:"# My Post\n\nThis is **markdown** content"`
	ContentJSON     string            `gorm:"column:content_json;type:text" json:"content_json" example:"{\"type\":\"doc\",\"content\":[]}"`
	Status          PostVersionStatus `gorm:"not null" json:"status" example:"draft"`
	AuthorID        uint              `gorm:"not null" json:"author_id" example:"1"`
	Post            *Post             `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Author          *User             `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	CreatedAt       time.Time         `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt       time.Time         `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// GetID implements the ModelInterface
func (pv PostVersion) GetID() uint {
	return pv.ID
}
