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

	logger.Debug("bcrypt.Hash: Success hash", "hashedLength", len(bytes), "hashed", string(bytes))
	return string(bytes), nil
}

func Compare(password string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	if err != nil {
		logger.Error("bcrypt.Compare: Error comparing hash and password", "error", err)
		logger.Debug("bcrypt.Compare: Error comparing hash and password", "hashedLength", len(hashed), "hashed", hashed)
	}

	return err == nil
}