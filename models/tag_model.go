package models

import "time"

type Tag struct {
	ID          uint      `gorm:"primaryKey" json:"id" example:"1"`
	Name        string    `gorm:"not null;unique" json:"name" example:"technology"`
	Description string    `gorm:"type:text" json:"description" example:"Technology related posts"`
	UsageCount  int       `gorm:"default:0" json:"usage_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostTag struct {
	PostID    uint      `gorm:"primaryKey" json:"post_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}
