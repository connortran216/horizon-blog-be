package test

import (
	"go-crud/initializers"
	"go-crud/models"
	"go-crud/services"

	"github.com/brianvoe/gofakeit/v6"
)

type PostOption func(*models.Post)

func WithTitle(title string) PostOption {
	return func(p *models.Post) {
		p.Title = title
	}
}

func WithContent(content string) PostOption {
	return func(p *models.Post) {
		p.Content = content
	}
}

func WithUserID(userID uint) PostOption {
	return func(p *models.Post) {
		p.UserID = userID
	}
}

func PostFactory(opts ...PostOption) models.Post {
	// Create a default user only if no UserID is explicitly set
	post := &models.Post{
		Title:   gofakeit.Sentence(6),
		Content: gofakeit.Paragraph(1, 3, 12, " "),
	}

	// Check if UserID is set via options
	userIDProvided := false
	for _, opt := range opts {
		opt(post)
		if post.UserID != 0 {
			userIDProvided = true
		}
	}

	// If no UserID was provided, create a default user and assign it
	if !userIDProvided {
		user := UserFactory("testpassword123")
		post.UserID = user.ID
	}

	result := initializers.DB.Create(post)
	if result.Error != nil {
		panic(result.Error) // Fail fast in tests
	}
	return *post
}

type UserOption func(*models.User)

func WithEmail(email string) UserOption {
	return func(u *models.User) {
		u.Email = email
	}
}

func WithName(name string) UserOption {
	return func(u *models.User) {
		u.Name = name
	}
}

func UserFactory(plainPassword string, opts ...UserOption) models.User {
	user := &models.User{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		HashedPassword: plainPassword, // Set plain password
	}

	for _, opt := range opts {
		opt(user)
	}

	// Hash the password before saving (same as service.Create)
	if hashedPassword, err := services.HashPassword(user.HashedPassword); err == nil {
		user.HashedPassword = hashedPassword
	}

	initializers.DB.Create(user)
	return *user
}
