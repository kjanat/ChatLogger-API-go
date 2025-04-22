package service

import (
	"errors"
	"fmt"
	"time"

	"ChatLogger-API-go/internal/domain"
	"ChatLogger-API-go/internal/hash"
	"github.com/golang-jwt/jwt/v5"
)

// UserService implements the domain.UserService interface.
type UserService struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

// NewUserService creates a new user service.
func NewUserService(userRepo domain.UserRepository, jwtSecret string) domain.UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register registers a new user.
func (s *UserService) Register(user *domain.User, password string) error {
	// Check if user with this email already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Hash the password using our centralized hash package
	hashedPassword, err := hash.GeneratePasswordHash(password, 10)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Set the hashed password and create timestamps
	user.PasswordHash = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Create the user
	return s.userRepo.Create(user)
}

// Login authenticates a user and returns a JWT token.
func (s *UserService) Login(email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return "", errors.New("invalid email or password")
	}

	// Check password using our centralized hash package
	if err := hash.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := generateJWT(user, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return token, nil
}

// GetByID gets a user by ID.
func (s *UserService) GetByID(id uint) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

// GetByEmail gets a user by email.
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(email)
}

// GetByOrganizationID gets users by organization ID with pagination.
func (s *UserService) GetByOrganizationID(orgID uint, limit, offset int) ([]domain.User, error) {
	return s.userRepo.FindByOrganizationID(orgID, limit, offset)
}

// UpdateUser updates a user.
func (s *UserService) UpdateUser(user *domain.User) error {
	// Set updated time
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// ChangePassword changes a user's password.
func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Get the user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Verify current password using our centralized hash package
	if err := hash.VerifyPassword(user.PasswordHash, currentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash the new password using our centralized hash package
	hashedPassword, err := hash.GeneratePasswordHash(newPassword, 10)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Update the user's password
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// DeleteUser deletes a user.
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

// JWTClaims represents the claims in a JWT token.
type JWTClaims struct {
	UserID         uint        `json:"uid"`
	Email          string      `json:"email"`
	OrganizationID uint        `json:"org_id"`
	Role           domain.Role `json:"role"`
	jwt.RegisteredClaims
}

// generateJWT generates a JWT token for a user.
func generateJWT(user *domain.User, secret string) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &JWTClaims{
		UserID:         user.ID,
		Email:          user.Email,
		OrganizationID: user.OrganizationID,
		Role:           user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
