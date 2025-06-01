package jwt

import (
	"errors"
	"fmt"
	"javaneseivankov/url-shortener/internal/errx"
	repository "javaneseivankov/url-shortener/internal/repository/model"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

type JWT struct {
	Secret []byte
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

	log.Println("TTL STRING IS " + ttlString)
	log.Println("TTL IS " + ttl.String())

	return JWT{
		Secret: []byte(secretKey),
		TTL: ttl,
	}
}

func (j *JWT) GenerateToken(user repository.User) (string, error) {
	log.Println("JWT TTL in Generate Token: " + j.TTL.String())

	claims := Claims{
		UserID: user.ID.String(),
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        user.ID.String(),
		},
	}
	claims.UserID = user.ID.String()
	claims.Email = user.Email
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(j.TTL))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secret)
}

func (j *JWT) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.Secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Print("Token expired.")
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			log.Print("Malformed token.")
		} else {
			log.Printf("JWT parse error: %v", err)
		}
		return nil, err
	}


	if !token.Valid {
		log.Print("Invalid bearer token: " + token.Raw)
		return nil, errx.ErrInvalidBearerToken
	}

	return claims, nil
}
