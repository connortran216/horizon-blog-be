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

// Create creates a new post metadata
func (s *PostService) Create(post models.Post, tagNames []string) (*models.Post, error) {

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Create the post (metadata only)
	result := tx.Create(&post)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Associate tags if provided
	if len(tagNames) > 0 {
		tagService := NewTagService()
		tagService.db = tx
		err := tagService.AssociateTagsWithPost(post.ID, tagNames)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tagService.db = s.db // Reset
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load the post with tags
	result = s.db.Preload("Tags").First(&post, post.ID)
	return &post, result.Error
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

// GetWithPagination retrieves posts with pagination
func (s *PostService) GetWithPagination(query schemas.ListPostsQueryParams) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	db := s.db.Model(&models.Post{})

	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}

	if len(query.TagNames) > 0 {
		tagService := NewTagService()
		var tagIDs []uint
		for _, tagName := range query.TagNames {
			tag, err := tagService.GetByName(tagName)
			if err == nil {
				tagIDs = append(tagIDs, tag.ID)
			}
		}

		if len(tagIDs) > 0 {
			db = db.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
				Where("post_tags.tag_id IN ?", tagIDs).
				Group("posts.id")
		}
	}

	if query.Status != nil {
		// Filter by post status based on whether versions exist with published status
		if *query.Status == models.Published {
			// Posts that have at least one published version
			db = db.Where("EXISTS (SELECT 1 FROM post_versions pv WHERE pv.post_id = posts.id AND pv.status = ?)", models.VersionPublished)
		} else if *query.Status == models.Draft {
			// Posts that don't have any published versions (draft-only posts)
			db = db.Where("NOT EXISTS (SELECT 1 FROM post_versions pv WHERE pv.post_id = posts.id AND pv.status = ?)", models.VersionPublished)
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Limit
	result := db.Preload("User").Preload("Tags").Limit(query.Limit).Offset(offset).Find(&posts)
	return posts, total, result.Error
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

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete associated post versions first (to avoid FK constraint violations)
	if err := tx.Unscoped().Where("post_id = ?", id).Delete(&models.PostVersion{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the post
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}
