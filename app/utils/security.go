package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type IPasswrodUtils interface {
	GenerateSalt(length int) (string, error)
	HashPassword(password string, salt string) (string, error)
	VerifyPassword(hashedPassword string, password string, salt string) bool
}

type PasswordUtils struct {
}

func (pu *PasswordUtils) GenerateSalt(length int) (string, error) {
	placeBytes := make([]byte, length)

	if _, err := rand.Read(placeBytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(placeBytes), nil
}

func (pu *PasswordUtils) HashPassword(password string, salt string) (string, error) {
	salted := password + salt
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(salted), bcrypt.DefaultCost)

	return string(hashBytes), err
}

func (pu *PasswordUtils) VerifyPassword(hashedPassword string, password string, salt string) bool {
	salted := password + salt

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(salted)); err != nil {
		return false
	}
	return true
}
