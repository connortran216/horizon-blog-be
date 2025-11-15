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

func WithContentMarkdown(content string) PostOption {
	return func(p *models.Post) {
		p.ContentMarkdown = content
	}
}

func WithContent(content string) PostOption {
	return WithContentMarkdown(content)
}

func WithContentJSON(content string) PostOption {
	return func(p *models.Post) {
		p.ContentJSON = content
	}
}

func WithUserID(userID uint) PostOption {
	return func(p *models.Post) {
		p.UserID = userID
	}
}

func WithStatus(status models.PostStatus) PostOption {
	return func(p *models.Post) {
		p.Status = status
	}
}

func PostFactory(opts ...PostOption) models.Post {
	// Create a default user only if no UserID is explicitly set
	post := &models.Post{
		Title:           gofakeit.Sentence(6),
		ContentMarkdown: gofakeit.Paragraph(1, 3, 12, " "),
		Status:          models.Published,
	}

	// Generate ContentJSON in Milkdown ProseMirror format
	docContent := []map[string]interface{}{}
	// Heading
	docContent = append(docContent, map[string]interface{}{
		"type":  "heading",
		"attrs": map[string]interface{}{"level": 1},
		"content": []map[string]interface{}{
			{"type": "text", "text": gofakeit.Sentence(4)},
		},
	})
	// Paragraph with marked text
	docContent = append(docContent, map[string]interface{}{
		"type": "paragraph",
		"content": []map[string]interface{}{
			{"type": "text", "text": "This is "},
			{
				"type":  "text",
				"marks": []map[string]interface{}{{"type": "strong"}},
				"text":  gofakeit.Word(),
			},
			{"type": "text", "text": " editor with "},
			{
				"type":  "text",
				"marks": []map[string]interface{}{{"type": "em"}},
				"text":  gofakeit.Word(),
			},
			{"type": "text", "text": " syntax."},
		},
	})
	// Bullet list
	docContent = append(docContent, map[string]interface{}{
		"type": "bullet_list",
		"content": []map[string]interface{}{
			{
				"type": "list_item",
				"content": []map[string]interface{}{
					{
						"type": "paragraph",
						"content": []map[string]interface{}{
							{"type": "text", "text": gofakeit.Word()},
						},
					},
				},
			},
			{
				"type": "list_item",
				"content": []map[string]interface{}{
					{
						"type": "paragraph",
						"content": []map[string]interface{}{
							{"type": "text", "text": gofakeit.Word()},
						},
					},
				},
			},
		},
	})
	doc := map[string]interface{}{
		"type":    "doc",
		"content": docContent,
	}
	jsonBytes, _ := json.Marshal(doc)
	post.ContentJSON = string(jsonBytes)

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
		Name:           gofakeit.Name(),
		Email:          gofakeit.Email(),
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
