package pkg

import (
	"fmt"
	"javaneseivankov/url-shortener/internal/errx"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateUserID() string {
	return ""
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email   string `json:"email"`
	Username string `json:"username"`
	Role string `json:"role"` 
	jwt.RegisteredClaims
}

type JWT struct {
	Secret string
	TTL time.Duration
}

func NewJWT(secretKey string, ttlString string) JWT {
	if secretKey == "" {
		log.Fatalln("Secret Key Cannot be Empty!")
	}

	ttl, err := time.ParseDuration(ttlString)
	if err != nil {
		log.Fatalln("Failed to parse TTL:", err)
	}

	return JWT{
		Secret: secretKey,
		TTL: ttl,
	}
}

func (j *JWT) GenerateToken(username string, password string) (string, error) {
	claims := Claims{}
	claims.Email = "foo@gmail.com"
	claims.Username = username
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.TTL))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secret)
}

func (j *JWT) VerifyToken(tokenString string, claims *Claims) (error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.Secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errx.ErrInvalidBearerToken
	}

	return nil
}
