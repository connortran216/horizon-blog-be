package views

import (
	"fmt"
	"go-crud/models"
	"go-crud/schemas"
	"go-crud/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


type UserViews struct {
	service *services.UserService
}

func NewUserViews() *UserViews {
	return &UserViews{
		service: services.NewUserService(),
	}
}


// @Summary Create user
// @Tags users
// @Param user body schemas.CreateUserInput true "User data"
// @Success 201 {object} schemas.AuthResponse
// @Router /users [post]
func (v *UserViews) CreateUser(c *gin.Context) {
	var input schemas.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Invalid request data: %v", err),
		})
		return
	}



	result, err := v.service.Create(models.User{
		Name:         input.Name,
		Email:        input.Email,
		HashedPassword: input.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to create user: %v", err),
		})
		return
	}

	// Generate JWT token for newly registered user
	token, err := services.GenerateToken(*result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to generate token: %v", err),
		})
		return
	}

	response := schemas.AuthResponse{
		Token: token,
		UserResponse: schemas.UserResponse{
			Data:    *result,
			Message: "User created successfully",
		},
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary Get user by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} schemas.UserResponse
// @Router /users/{id} [get]
func (v *UserViews) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	result, err := v.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: fmt.Sprintf("User not found: %v", err),
		})
		return
	}

	response := schemas.UserResponse{
		Data:    *result,
		Message: "User retrieved successfully",
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Partially update user
// @Tags users
// @Param id path int true "User ID"
// @Param user body schemas.PartialUpdateUserInput true "User data"
// @Success 200 {object} schemas.UserResponse
// @Router /users/{id} [patch]
func (v *UserViews) PartialUpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Check if target user exists first
	_, err = v.service.GetByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, schemas.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to retrieve user: %v", err),
		})
		return
	}

	// Then check if user is authorized to update this account
	authenticatedUserID, exists := GetUserIDFromContext(c)
	if !exists || authenticatedUserID != uint(id) {
		c.JSON(http.StatusForbidden, schemas.ErrorResponse{
			Error: "You can only update your own account",
		})
		return
	}

	var input schemas.PartialUpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: fmt.Sprintf("Invalid request data: %v", err),
		})
		return
	}



	result, err := v.service.PartialUpdate(uint(id), input)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, schemas.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to update user: %v", err),
		})
		return
	}

	response := schemas.UserResponse{
		Data:    *result,
		Message: "User updated successfully",
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Delete user
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} schemas.MessageResponse
// @Router /users/{id} [delete]
func (v *UserViews) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Check if target user exists first
	_, err = v.service.GetByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, schemas.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to retrieve user: %v", err),
		})
		return
	}

	// Then check if user is authorized to delete this account
	authenticatedUserID, exists := GetUserIDFromContext(c)
	if !exists || authenticatedUserID != uint(id) {
		c.JSON(http.StatusForbidden, schemas.ErrorResponse{
			Error: "You can only delete your own account",
		})
		return
	}

	if err := v.service.Delete(uint(id)); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, schemas.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		c.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to delete user: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.MessageResponse{
		Message: "User deleted successfully",
	})
}

// @Summary List user's posts by status
// @Tags users
// @Param status query string false "Post status (draft or published)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} schemas.ListPostsResponse
// @Router /users/me/posts [get]
func (v *UserViews) ListUserPosts(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	var query schemas.ListPostsQueryParams
	query.Page, _ = strconv.Atoi(c.Query("page"))
	query.Limit, _ = strconv.Atoi(c.Query("limit"))

	// Set defaults if not provided
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	// Set user ID to authenticated user
	query.UserID = &userID

	// Handle status filtering (Note: status filtering not yet implemented in current PostService.GetWithPagination)
	// TODO: Implement status-based post filtering if needed
	// if statusParam := c.Query("status"); statusParam != "" {
	//     status := models.PostStatus(statusParam)
	//     // Validate status value
	//     if status != models.Draft && status != models.Published {
	//         c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
	//             Error: "Invalid status: must be 'draft' or 'published'",
	//         })
	//         return
	//     }
	//     query.Status = &status
	// }

	// Use post service to get posts
	postService := services.NewPostService()
	results, total, err := postService.GetWithPagination(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: fmt.Sprintf("Failed to fetch posts: %v", err),
		})
		return
	}

	response := schemas.ListPostsResponse{
		Data:  results,
		Limit: query.Limit,
		Page:  query.Page,
		Total: int(total),
	}
	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers user-related routes
func (v *UserViews) RegisterRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.POST("", v.CreateUser)
		users.GET("/me/posts", AuthMiddleware(), v.ListUserPosts)
		users.GET("/:id", v.GetUserByID)
		users.PATCH("/:id", AuthMiddleware(), v.PartialUpdateUser)
		users.DELETE("/:id", AuthMiddleware(), v.DeleteUser)
	}
}
