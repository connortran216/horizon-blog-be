package models

import "time"

type Post struct {
	ID          uint      `gorm:"primaryKey" json:"id" example:"1"`
	UserID      uint      `gorm:"not null" json:"user_id" example:"1"`
	Title       string    `gorm:"not null" json:"title" example:"My First Post"`
	Content     string    `gorm:"not null" json:"content" example:"This is the content of my first post"`
	IsDraft     bool      `gorm:"default:true;not null" json:"is_draft" example:"false"`
	IsPublished bool      `gorm:"default:false;not null" json:"is_published" example:"true"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tags        []Tag     `gorm:"many2many:post_tags" json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// GetID implements the ModelInterface
func (p Post) GetID() uint {
	return p.ID
}
