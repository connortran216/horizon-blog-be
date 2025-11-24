package services

import (
	"errors"
	"go-crud/initializers"
	"go-crud/models"
	"go-crud/schemas"

	"gorm.io/gorm"
)

// PostService handles business logic for Post operations
type PostService struct {
	db *gorm.DB
}

// NewPostService creates a new PostService instance
func NewPostService() *PostService {
	return &PostService{
		db: initializers.DB,
	}
}

// Create creates a new post
func (s *PostService) Create(post models.Post, tagNames []string) (*models.Post, error) {
	if post.Title == "" {
		return nil, errors.New("title is required")
	}
	if post.ContentMarkdown == "" {
		return nil, errors.New("content_markdown is required")
	}
	if post.ContentJSON == "" {
		return nil, errors.New("content_json is required")
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Create the post
	result := tx.Create(&post)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Associate tags if provided
	if len(tagNames) > 0 {
		tagService := NewTagService()
		err := tagService.AssociateTagsWithPost(post.ID, tagNames)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load the post with tags
	return s.GetByID(post.ID)
}

// GetByID retrieves a post by ID
func (s *PostService) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	result := s.db.Preload("User").Preload("Tags").First(&post, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, result.Error
	}

	return &post, nil
}

// GetAll retrieves all posts
func (s *PostService) GetAll() ([]models.Post, error) {
	var posts []models.Post
	result := s.db.Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

// GetPaginated retrieves posts with pagination
func (s *PostService) GetWithPagination(query schemas.ListPostsQueryParams) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Build query with optional filters
	db := s.db.Model(&models.Post{})

	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	// Filter by tags if provided
	if len(query.TagNames) > 0 {
		// Get tag IDs from tag names
		tagService := NewTagService()
		var tagIDs []uint
		for _, tagName := range query.TagNames {
			tag, err := tagService.GetByName(tagName)
			if err == nil {
				tagIDs = append(tagIDs, tag.ID)
			}
		}

		if len(tagIDs) > 0 {
			// Join with post_tags and filter by tag IDs
			db = db.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
				Where("post_tags.tag_id IN ?", tagIDs).
				Group("posts.id")
		}
	}

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (query.Page - 1) * query.Limit

	// Get paginated results with preloading, sorted by created date DESC
	result := db.Preload("User").Preload("Tags").Order("created_at DESC").Limit(query.Limit).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return posts, total, nil
}

// Update updates an existing post
func (s *PostService) Update(id uint, updatedPost models.Post) (*models.Post, error) {
	var post models.Post

	// Check if post exists
	result := s.db.First(&post, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, result.Error
	}

	// Validate updated data
	if updatedPost.Title == "" {
		return nil, errors.New("title is required")
	}
	if updatedPost.ContentMarkdown == "" {
		return nil, errors.New("content_markdown is required")
	}
	if updatedPost.ContentJSON == "" {
		return nil, errors.New("content_json is required")
	}

	// Update fields
	post.Title = updatedPost.Title
	post.ContentMarkdown = updatedPost.ContentMarkdown
	post.ContentJSON = updatedPost.ContentJSON
	if updatedPost.Status != "" {
		post.Status = updatedPost.Status
	}

	// Save changes
	result = s.db.Save(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	return &post, nil
}

// PartialUpdate updates specific fields of an existing post
func (s *PostService) PartialUpdate(id uint, partialData map[string]interface{}) (*models.Post, error) {
	var post models.Post

	// Check if post exists
	result := s.db.First(&post, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, result.Error
	}

	// Update only provided fields
	if title, exists := partialData["title"]; exists {
		if titleStr, ok := title.(string); ok && titleStr != "" {
			post.Title = titleStr
		} else if titleStr == "" {
			return nil, errors.New("title cannot be empty")
		}
	}

	if contentMarkdown, exists := partialData["content_markdown"]; exists {
		if contentStr, ok := contentMarkdown.(string); ok && contentStr != "" {
			post.ContentMarkdown = contentStr
		} else if contentStr == "" {
			return nil, errors.New("content_markdown cannot be empty")
		}
	}

	if contentJSON, exists := partialData["content_json"]; exists {
		if contentStr, ok := contentJSON.(string); ok && contentStr != "" {
			post.ContentJSON = contentStr
		} else if contentStr == "" {
			return nil, errors.New("content_json cannot be empty")
		}
	}

	if status, exists := partialData["status"]; exists {
		if statusEnum, ok := status.(models.PostStatus); ok && (statusEnum == models.Draft || statusEnum == models.Published) {
			post.Status = statusEnum
		} else {
			return nil, errors.New("invalid status: must be 'draft' or 'published'")
		}
	}

	// Save changes
	result = s.db.Save(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	return &post, nil
}

// Delete deletes a post by ID
func (s *PostService) Delete(id uint) error {
	var post models.Post

	// Check if post exists
	result := s.db.First(&post, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return result.Error
	}

	// Delete the post
	result = s.db.Delete(&post)
	return result.Error
}
