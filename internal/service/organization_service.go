package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"ChatLogger-API-go/internal/domain"
)

// OrganizationService implements the domain.OrganizationService interface.
type OrganizationService struct {
	orgRepo domain.OrganizationRepository
}

// NewOrganizationService creates a new organization service.
func NewOrganizationService(orgRepo domain.OrganizationRepository) domain.OrganizationService {
	return &OrganizationService{
		orgRepo: orgRepo,
	}
}

// Create creates a new organization.
func (s *OrganizationService) Create(org *domain.Organization) error {
	// Generate a slug if not provided
	if org.Slug == "" {
		org.Slug = generateSlug(org.Name)
	}

	// Check if organization with this slug already exists
	existingOrg, err := s.orgRepo.FindBySlug(org.Slug)
	if err != nil {
		return fmt.Errorf("error checking existing organization: %w", err)
	}

	if existingOrg != nil {
		return errors.New("organization with this slug already exists")
	}

	// Set timestamps
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	// Create the organization
	return s.orgRepo.Create(org)
}

// GetByID gets an organization by ID.
func (s *OrganizationService) GetByID(id uint) (*domain.Organization, error) {
	return s.orgRepo.FindByID(id)
}

// GetBySlug gets an organization by slug.
func (s *OrganizationService) GetBySlug(slug string) (*domain.Organization, error) {
	return s.orgRepo.FindBySlug(slug)
}

// Update updates an organization.
func (s *OrganizationService) Update(org *domain.Organization) error {
	// Get the existing organization
	existingOrg, err := s.orgRepo.FindByID(org.ID)
	if err != nil {
		return fmt.Errorf("error finding organization: %w", err)
	}

	if existingOrg == nil {
		return errors.New("organization not found")
	}

	// If slug is being changed, check if the new slug already exists
	if org.Slug != "" && org.Slug != existingOrg.Slug {
		orgWithSlug, err := s.orgRepo.FindBySlug(org.Slug)
		if err != nil {
			return fmt.Errorf("error checking slug: %w", err)
		}

		if orgWithSlug != nil {
			return errors.New("organization with this slug already exists")
		}
	}

	// Update timestamp
	org.UpdatedAt = time.Now()

	// Update the organization
	return s.orgRepo.Update(org)
}

// Delete deletes an organization.
func (s *OrganizationService) Delete(id uint) error {
	return s.orgRepo.Delete(id)
}

// List lists organizations with pagination.
func (s *OrganizationService) List(limit, offset int) ([]domain.Organization, error) {
	return s.orgRepo.List(limit, offset)
}

// Helper functions

// generateSlug generates a URL-friendly slug from a name.
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with dashes
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	reg, _ := regexp.Compile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple dashes
	reg, _ = regexp.Compile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim dashes from beginning and end
	slug = strings.Trim(slug, "-")

	return slug
}
