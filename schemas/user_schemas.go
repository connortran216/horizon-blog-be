package schemas

import "go-crud/models"

type CreateUserInput struct {
	Name         string `json:"name" binding:"required" example:"Connor Tran"`
	Email        string `json:"email" binding:"required,email" example:"connor@example.com"`
	Password string `json:"password" binding:"required" example:"abcxyz123"`
}

type PartialUpdateUserInput struct {
	Name  *string `json:"name" binding:"omitempty,min=3" example:"Connor Tran"`
	Email *string `json:"email" binding:"omitempty,email" example:"connor@example.com"`
	Password *string `json:"password" example:"abcxyz123"`
}

type UserResponse struct {
	Data    models.User `json:"data"`
	Message string      `json:"message" example:"User created successfully"`
}
