package bcrypt

import (
	"javaneseivankov/url-shortener/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func Hash(plainText string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)

	if err != nil {
		logger.Error("bcrypt.Hash: Failed to hash password", "error", err)
		return "", err
	}

	return string(bytes), nil
}

func Compare(password string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}