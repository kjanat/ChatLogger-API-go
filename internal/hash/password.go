// Package hash provides password hashing and verification functionality
// for the ChatLogger API using bcrypt for secure password storage.
package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash creates a bcrypt hash from a plain-text password.
func GeneratePasswordHash(password string, cost int) (string, error) {
	if cost < 4 {
		cost = 10 // Default to cost 10 if provided cost is too low
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword checks if the provided password matches the stored hash.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
