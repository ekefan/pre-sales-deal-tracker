package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword uses bcrypt to create a hash for password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPasswordHash uses bcrypt to compare password with hash
// If compatible true is returned
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
