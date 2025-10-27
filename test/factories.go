package test

import (
	"encoding/json"
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



func WithUserID(userID uint) PostOption {
	return func(p *models.Post) {
		p.UserID = userID
	}
}

func PostFactory(opts ...PostOption) models.Post {
	post := &models.Post{
		Title: gofakeit.Sentence(6),
	}

	// Generate Content for versioning (keep backward compatibility in tests)
	contentMarkdown := gofakeit.Paragraph(1, 3, 12, " ")
	// Generate ContentJSON in Milkdown ProseMirror format
	docContent := []map[string]interface{}{}
	docContent = append(docContent, map[string]interface{}{
		"type": "heading",
		"attrs": map[string]interface{}{"level": 1},
		"content": []map[string]interface{}{
			{"type": "text", "text": gofakeit.Sentence(4)},
		},
	})
	docContent = append(docContent, map[string]interface{}{
		"type": "paragraph",
		"content": []map[string]interface{}{
			{"type": "text", "text": gofakeit.Sentence(10)},
		},
	})
	doc := map[string]interface{}{
		"type": "doc",
		"content": docContent,
	}
	jsonBytes, _ := json.Marshal(doc)
	contentJSON := string(jsonBytes)

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

	// Create post metadata first
	result := initializers.DB.Create(post)
	if result.Error != nil {
		panic(result.Error)
	}

	// Create initial version for the post
	versionService := services.NewPostVersionService()
	version, err := versionService.CreateDraftVersion(post.ID, post.UserID, post.Title, contentMarkdown, contentJSON)
	if err != nil {
		panic(err)
	}

	// Update the title from the version for consistency
	post.Title = version.Title

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
