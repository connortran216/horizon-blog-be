package services

import (
	"errors"
	"go-crud/initializers"
	"go-crud/models"
	"go-crud/schemas"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		db: initializers.DB,
	}
}

// Create creates a new user
func (s *UserService) Create(user models.User) (*models.User, error) {
	if user.Name == "" {
		return nil, errors.New("name is required")
	}
	if user.Email == "" {
		return nil, errors.New("email is required")
	}
	if user.HashedPassword == "" {
		return nil, errors.New("password is required")
	}

	hashedPassword, err := HashPassword(user.HashedPassword)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user.HashedPassword = hashedPassword

	return &user, s.db.Create(&user).Error
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Partial Update user fields partially
func (s *UserService) PartialUpdate(id uint, input schemas.PartialUpdateUserInput) (*models.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		hashedPassword, err := HashPassword(*input.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.HashedPassword = hashedPassword
	}

	result := s.db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// Delete user by ID
func (s *UserService) Delete(id uint) error {
	var user models.User

	result := s.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return result.Error
	}

	result = s.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByEmail finds a user by email address
func (s *UserService) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}
