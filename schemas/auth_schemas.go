package schemas

type LoginInput struct {
	Email    string `json:"email" validate:"required,email" example:"connor@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type AuthResponse struct {
	Token string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}
