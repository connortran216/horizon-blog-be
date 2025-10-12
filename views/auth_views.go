package views

import (
	"fmt"
	"go-crud/schemas"
	"go-crud/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthViews struct {
	userService *services.UserService
	validator   *validator.Validate
}

func NewAuthViews() *AuthViews {
	return &AuthViews{
		userService: services.NewUserService(),
		validator:   validator.New(),
	}
}

// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param loginInput body schemas.LoginInput true "Login credentials"
// @Success 200 {object} schemas.AuthResponse
// @Router /auth/login [post]
func (v *AuthViews) Login(c *gin.Context) {
	var input schemas.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Invalid request data: %v", err),
		})
		return
	}

	if err := v.validator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Validation failed: %v", err),
		})
		return
	}

	// Find user by email - need to modify UserService to add FindByEmail method
	user, err := v.userService.FindByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "Invalid email or password",
		})
		return
	}

	// Check password
	if !services.CheckHashedPassword(input.Password, user.HashedPassword) {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := services.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to generate token: %v", err),
		})
		return
	}

	response := schemas.AuthResponse{
		Token: token,
		User: schemas.UserResponse{
			Data:    *user,
			Message: "Login successful",
		},
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers auth-related routes
func (v *AuthViews) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", v.Login)
	}
}
