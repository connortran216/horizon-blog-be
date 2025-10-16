package services

import (
	"errors"
	"go-crud/initializers"
	"go-crud/models"
	"go-crud/schemas"
	"strings"

	"gorm.io/gorm"
)

// TagService handles business logic for Tag operations
type TagService struct {
	db *gorm.DB
}

// NewTagService creates a new TagService instance
func NewTagService() *TagService {
	return &TagService{
		db: initializers.DB,
	}
}

// Create creates a new tag
func (s *TagService) Create(tag models.Tag) (*models.Tag, error) {
	// Validate tag name
	if tag.Name == "" {
		return nil, errors.New("tag name is required")
	}

	// Normalize tag name (lowercase, trim spaces)
	tag.Name = strings.ToLower(strings.TrimSpace(tag.Name))

	// Check if tag already exists
	var existingTag models.Tag
	result := s.db.Where("name = ?", tag.Name).First(&existingTag)
	if result.Error == nil {
		return nil, errors.New("tag already exists")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// Create the tag
	result = s.db.Create(&tag)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tag, nil
}

// GetByID retrieves a tag by ID
func (s *TagService) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	result := s.db.First(&tag, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, result.Error
	}

	return &tag, nil
}

// GetByName retrieves a tag by name
func (s *TagService) GetByName(name string) (*models.Tag, error) {
	var tag models.Tag
	result := s.db.Where("name = ?", strings.ToLower(strings.TrimSpace(name))).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, result.Error
	}

	return &tag, nil
}

// GetAll retrieves all tags with pagination and sorting
func (s *TagService) GetAll(query schemas.TagListQueryParams) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64

	// Build query
	db := s.db.Model(&models.Tag{})

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	switch query.Sort {
	case "usage_count":
		db = db.Order("usage_count DESC")
	case "created_at":
		db = db.Order("created_at DESC")
	default:
		db = db.Order("name ASC")
	}

	// Apply pagination
	offset := (query.Page - 1) * query.Limit
	result := db.Limit(query.Limit).Offset(offset).Find(&tags)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return tags, total, nil
}

// GetPopular retrieves most used tags
func (s *TagService) GetPopular(limit int) ([]models.Tag, error) {
	var tags []models.Tag
	result := s.db.Order("usage_count DESC").Limit(limit).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}

	return tags, nil
}

// Search searches tags by name
func (s *TagService) Search(query string, limit int) ([]models.Tag, error) {
	var tags []models.Tag
	searchQuery := "%" + strings.ToLower(query) + "%"
	result := s.db.Where("LOWER(name) LIKE ?", searchQuery).Limit(limit).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}

	return tags, nil
}

// IncrementUsage increments the usage count for a tag
func (s *TagService) IncrementUsage(tagID uint) error {
	result := s.db.Model(&models.Tag{}).Where("id = ?", tagID).Update("usage_count", gorm.Expr("usage_count + 1"))
	return result.Error
}

// Update updates an existing tag
func (s *TagService) Update(id uint, updates map[string]interface{}) (*models.Tag, error) {
	var tag models.Tag

	// Check if tag exists
	result := s.db.First(&tag, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, result.Error
	}

	// Normalize name if provided
	if name, exists := updates["name"]; exists {
		if nameStr, ok := name.(string); ok {
			updates["name"] = strings.ToLower(strings.TrimSpace(nameStr))
		}
	}

	// Update the tag
	result = s.db.Model(&tag).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return updated tag
	return s.GetByID(id)
}

// Delete deletes a tag by ID
func (s *TagService) Delete(id uint) error {
	// Check if tag exists
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Delete the tag (this will cascade to post_tags due to foreign key)
	result := s.db.Delete(&models.Tag{}, id)
	return result.Error
}

// GetOrCreateTags gets existing tags or creates new ones
func (s *TagService) GetOrCreateTags(tagNames []string) ([]models.Tag, error) {
	var tags []models.Tag

	for _, name := range tagNames {
		if name == "" {
			continue
		}

		normalizedName := strings.ToLower(strings.TrimSpace(name))

		// Try to find existing tag
		tag, err := s.GetByName(normalizedName)
		if err != nil {
			if err.Error() == "tag not found" {
				// Create new tag
				newTag := models.Tag{Name: normalizedName}
				createdTag, err := s.Create(newTag)
				if err != nil {
					return nil, err
				}
				tags = append(tags, *createdTag)
			} else {
				return nil, err
			}
		} else {
			tags = append(tags, *tag)
		}
	}

	return tags, nil
}

// AssociateTagsWithPost associates tags with a post
func (s *TagService) AssociateTagsWithPost(postID uint, tagNames []string) error {
	// Get or create tags
	tags, err := s.GetOrCreateTags(tagNames)
	if err != nil {
		return err
	}

	// Remove existing associations
	err = s.db.Where("post_id = ?", postID).Delete(&models.PostTag{}).Error
	if err != nil {
		return err
	}

	// Create new associations
	for _, tag := range tags {
		postTag := models.PostTag{
			PostID: postID,
			TagID:  tag.ID,
		}
		err = s.db.Create(&postTag).Error
		if err != nil {
			return err
		}

		// Increment usage count
		s.IncrementUsage(tag.ID)
	}

	return nil
}
