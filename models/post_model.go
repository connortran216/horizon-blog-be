package models

import "time"

type PostStatus string

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)

type Post struct {
	ID                   uint         `gorm:"primaryKey" json:"id" example:"1"`
	UserID               uint         `gorm:"not null" json:"user_id" example:"1"`
	Title                string       `gorm:"not null" json:"title" example:"My First Post"`
	Slug                 string       `json:"slug,omitempty" example:"my-first-post"`
	PublishedVersionID   *uint        `json:"published_version_id,omitempty" example:"1"`
	User                 *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tags                 []Tag        `gorm:"many2many:post_tags" json:"tags,omitempty"`
	Versions             []PostVersion `gorm:"foreignKey:PostID" json:"versions,omitempty"`
	PublishedVersion     *PostVersion `gorm:"foreignKey:PublishedVersionID" json:"published_version,omitempty"`
	CreatedAt            time.Time    `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt            time.Time    `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// GetID implements the ModelInterface
func (p Post) GetID() uint {
	return p.ID
}
