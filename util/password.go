package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}