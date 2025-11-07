package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id" example:"1"`
	Name           string    `gorm:"not null" json:"name" example:"Connor Tran"`
	Email          string    `gorm:"unique;not null" json:"email" example:"connortran@gmail.com"`
	HashedPassword string    `gorm:"not null" json:"-"`
	CreatedAt      time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func (u User) GetID() uint {
	return u.ID
}
