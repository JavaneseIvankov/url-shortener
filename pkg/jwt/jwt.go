package jwt

import (
	"errors"
	"fmt"
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

type JWT struct {
    Secret []byte
    TTL    time.Duration
}

func NewJWT(secretKey string, ttlString string) JWT {
    if secretKey == "" {
        logger.Error("Secret Key Cannot be Empty!")
        panic("Secret Key Cannot be Empty!")
    }

    ttl, err := time.ParseDuration(ttlString)
    if err != nil {
        logger.Error("Failed to parse TTL", "error", err)
        panic(fmt.Sprintf("Failed to parse TTL: %v", err))
    }

    logger.Info("JWT.NewJWT: JWT initialized", "TTLString", ttlString, "TTL", ttl.String())

    return JWT{
        Secret: []byte(secretKey),
        TTL:    ttl,
    }
}

func (j *JWT) GenerateToken(user *model.User) (string, error) {
    logger.Debug("Generating token", "TTL", j.TTL.String())

    claims := Claims{
        UserID: user.ID.String(),
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TTL)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            ID:        user.ID.String(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(j.Secret)
    if err != nil {
        logger.Error("Failed to sign token", "error", err)
        return "", err
    }

    logger.Info("JWT.GenerateToken: Token generated successfully", "userID", user.ID.String())
    return signedToken, nil
}

func (j *JWT) VerifyToken(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            logger.Error("Unexpected signing method", "method", token.Header["alg"])
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return j.Secret, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            logger.Error("Token expired", "error", err)
            return nil, errx.ErrExpiredBearerToken
        } else if errors.Is(err, jwt.ErrTokenMalformed) {
            logger.Error("Token malformed", "error", err)
            return nil, errx.ErrxMalformedBearerToken
        } else {
            logger.Error("JWT parse error", "error", err)
        }
        return nil, err
    }

    if !token.Valid {
        logger.Warn("Invalid bearer token", "token", token.Raw)
        return nil, errx.ErrInvalidBearerToken
    }

    logger.Info("JWT.VerifyToken: Token verified successfully", "claims", claims)
    return claims, nil
}