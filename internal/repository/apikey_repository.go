package repository

import (
	"ChatLogger-API-go/internal/domain"
	"errors"
	"time"

	"gorm.io/gorm"
)

// APIKeyRepo implements the domain.APIKeyRepository interface
type APIKeyRepo struct {
	db *Database
}

// NewAPIKeyRepository creates a new API key repository
func NewAPIKeyRepository(db *Database) domain.APIKeyRepository {
	return &APIKeyRepo{db: db}
}

// Create creates a new API key
func (r *APIKeyRepo) Create(key *domain.APIKey) error {
	return r.db.Create(key).Error
}

// FindByID finds an API key by ID
func (r *APIKeyRepo) FindByID(id uint) (*domain.APIKey, error) {
	var key domain.APIKey
	err := r.db.First(&key, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &key, nil
}

// FindByHashedKey finds an API key by its hashed value
func (r *APIKeyRepo) FindByHashedKey(hashedKey string) (*domain.APIKey, error) {
	var key domain.APIKey
	err := r.db.Where("hashed_key = ?", hashedKey).First(&key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Don't return revoked keys
	if key.RevokedAt != nil {
		return nil, nil
	}

	return &key, nil
}

// ListByOrganizationID lists API keys for an organization
func (r *APIKeyRepo) ListByOrganizationID(orgID uint) ([]domain.APIKey, error) {
	var keys []domain.APIKey
	err := r.db.Where("organization_id = ?", orgID).Find(&keys).Error
	return keys, err
}

// Revoke revokes an API key by ID
func (r *APIKeyRepo) Revoke(id uint) error {
	now := time.Now()
	return r.db.Model(&domain.APIKey{}).Where("id = ?", id).Update("revoked_at", &now).Error
}

// Delete permanently deletes an API key by ID
func (r *APIKeyRepo) Delete(id uint) error {
	return r.db.Delete(&domain.APIKey{}, id).Error
}
