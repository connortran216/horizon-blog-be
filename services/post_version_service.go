package services

import (
	"errors"
	"go-crud/initializers"
	"go-crud/models"

	"gorm.io/gorm"
)

// PostVersionService handles business logic for PostVersion operations
type PostVersionService struct {
	db *gorm.DB
}

// NewPostVersionService creates a new PostVersionService instance
func NewPostVersionService() *PostVersionService {
	return &PostVersionService{
		db: initializers.DB,
	}
}

// CreateDraftVersion creates a new draft version for a post
// If no draft exists, creates a new one from the latest published version or empty
func (s *PostVersionService) CreateDraftVersion(postID, authorID uint, title, contentMarkdown, contentJSON string) (*models.PostVersion, error) {
	var post models.Post
	result := s.db.First(&post, postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, result.Error
	}

	// Check if a draft already exists for this post
	var existingDraft models.PostVersion
	draftResult := s.db.Where("post_id = ? AND status = ?", postID, models.VersionDraft).First(&existingDraft)

	if draftResult.Error == nil {
		// Return existing draft
		existingDraft.ContentMarkdown = contentMarkdown
		existingDraft.ContentJSON = contentJSON
		existingDraft.Title = title
		s.db.Save(&existingDraft)
		return &existingDraft, nil
	}

	// No existing draft, create new one
	var versionToCopy *models.PostVersion
	publishedResult := s.db.Where("post_id = ? AND status = ?", postID, models.VersionPublished).First(&versionToCopy)
	if publishedResult.Error != nil {
		// No published version found, set to nil
		versionToCopy = nil
	}

	var newVersion models.PostVersion
	newVersion.PostID = postID
	newVersion.AuthorID = authorID
	newVersion.Status = models.VersionDraft

	if versionToCopy != nil {
		newVersion.ContentMarkdown = contentMarkdown
		newVersion.ContentJSON = contentJSON
		newVersion.Title = title
	} else {
		// First version for this post
		newVersion.ContentMarkdown = contentMarkdown
		newVersion.ContentJSON = contentJSON
		newVersion.Title = title
	}

	createResult := s.db.Create(&newVersion)
	if createResult.Error != nil {
		return nil, createResult.Error
	}

	return &newVersion, nil
}

// AutoSaveDraft saves changes to an existing draft version
func (s *PostVersionService) AutoSaveDraft(versionID, authorID uint, title, contentMarkdown, contentJSON string) (*models.PostVersion, error) {
	var version models.PostVersion
	result := s.db.First(&version, versionID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("version not found")
		}
		return nil, result.Error
	}

	// Check ownership
	if version.AuthorID != authorID {
		return nil, errors.New("you can only update your own versions")
	}

	// Check if it's a draft
	if version.Status != models.VersionDraft {
		return nil, errors.New("can only auto-save draft versions")
	}

	// Update content
	version.Title = title
	version.ContentMarkdown = contentMarkdown
	version.ContentJSON = contentJSON

	result = s.db.Save(&version)
	return &version, result.Error
}

// PublishVersion publishes a version and marks all other versions for the same post as drafts
func (s *PostVersionService) PublishVersion(versionID, authorID uint) (*models.PostVersion, error) {
	var version models.PostVersion
	result := s.db.First(&version, versionID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("version not found")
		}
		return nil, result.Error
	}

	// Check ownership
	if version.AuthorID != authorID {
		return nil, errors.New("you can only publish your own versions")
	}

	// Begin transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Mark this version as published
	version.Status = models.VersionPublished
	if err := tx.Save(&version).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Mark all other versions for this post as drafts (only one can be published)
	if err := tx.Model(&models.PostVersion{}).
		Where("post_id = ? AND id != ?", version.PostID, versionID).
		Update("status", models.VersionDraft).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &version, nil
}

// GetByID retrieves a version by ID with related data
func (s *PostVersionService) GetByID(id uint) (*models.PostVersion, error) {
	var version models.PostVersion
	result := s.db.Preload("Post").Preload("Author").First(&version, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("version not found")
		}
		return nil, result.Error
	}

	return &version, nil
}

// GetVersionsForPost lists all versions for a post
func (s *PostVersionService) GetVersionsForPost(postID uint) ([]models.PostVersion, error) {
	var versions []models.PostVersion
	result := s.db.Where("post_id = ?", postID).
		Preload("Author").
		Order("created_at DESC").
		Find(&versions)

	if result.Error != nil {
		return nil, result.Error
	}

	return versions, nil
}

// GetWithPagination retrieves versions with pagination
func (s *PostVersionService) GetWithPagination(page, limit int, userID *uint, status *models.PostVersionStatus) ([]models.PostVersion, int64, error) {
	var versions []models.PostVersion
	var total int64

	db := s.db.Model(&models.PostVersion{})

	if userID != nil {
		db = db.Where("author_id = ?", *userID)
	}

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	result := db.Preload("Post").Preload("Author").Limit(limit).Offset(offset).Find(&versions)
	return versions, total, result.Error
}
