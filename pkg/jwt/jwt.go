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
    AccessTokenTTL    time.Duration
    RefreshTokenTTL    time.Duration
}

func NewJWT(secretKey string, accessTokenTTL string, refreshTokenTTL string) JWT {
    if secretKey == "" {
        logger.Error("jwt.NewJWT: Secret Key Cannot be Empty!")
        panic("Secret Key Cannot be Empty!")
    }

    aTTL, err := time.ParseDuration(accessTokenTTL)
    if err != nil {
        logger.Error("jwt.NewJWT: Failed to parse access token TTL", "error", err)
        panic(fmt.Sprintf("Failed to parse access token TTL: %v", err))
    }
	 logger.Debug("jwt.NewJWT: Parsed accessTokenTTL", "token", aTTL)

	 rTTL, err := time.ParseDuration(refreshTokenTTL)
	 if err != nil {
        logger.Error("jwt.NewJWT: Failed to parse refresh token TTL", "error", err)
        panic(fmt.Sprintf("Failed to parse refresh token TTL: %v", err))
	 }
	 logger.Debug("jwt.NewJWT: Parsed refreshTokenTTL", "token", rTTL)

    logger.Info("JWT.NewJWT: JWT initialized", "TTLString", accessTokenTTL, "TTL", aTTL.String())

    return JWT{
        Secret: []byte(secretKey),
        AccessTokenTTL:  aTTL,
		  RefreshTokenTTL: rTTL,
    }
}

func (j *JWT) createClaims(userID, email, id string, ttl time.Duration) *Claims {
    return &Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            ID:        id,
        },
    }
}

func (j *JWT) GenerateAccessToken(user *model.User) (string, error) {
    logger.Debug("jwt.GenerateAccessToken: Generating token", "TTL", j.AccessTokenTTL.String())

    claims := j.createClaims(user.ID.String(), user.Email, user.ID.String(), j.AccessTokenTTL)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(j.Secret)
    if err != nil {
        logger.Error("jwt.GenerateAccessToken: Failed to sign token", "error", err)
        return "", err
    }

    logger.Info("JWT.GenerateToken: Token generated successfully", "userID", user.ID.String())
    return signedToken, nil
}

func (j *JWT) GenerateRefreshToken(user *model.User) (string, error) {
    logger.Debug("jwt.GenerateRefreshToken: Generating token", "TTL", j.RefreshTokenTTL.String())

    claims := j.createClaims(user.ID.String(), user.Email, user.ID.String(), j.RefreshTokenTTL)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(j.Secret)
    if err != nil {
        logger.Error("jwt.GenerateRefreshToken: Failed to sign token", "error", err)
        return "", err
    }

    logger.Info("JWT.GenerateRefreshToken: Token generated successfully", "userID", user.ID.String())
    return signedToken, nil
}

func (j *JWT) RenewAccessToken(refreshToken string) (string, error) {
    logger.Debug("jwt.RenewAccessToken: Renewing access token")
    logger.Debug("jwt.RenewAccessToken: Renewing access token", "refreshToken", refreshToken)

	// FIXME: faulty logic, why verify already expired 
    claims, err := j.VerifyToken(refreshToken)
    if err != nil {
        return "", err
    }
	 claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.RefreshTokenTTL))


    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(j.Secret)
    if err != nil {
        logger.Error("jwt.RenewAccessToken: Failed to generate new access token", "error", err)
        return "", err
    }

    logger.Info("jwt.RenewAccessToken: Access token renewed successfully")
    return signedToken, nil
}


// This only needed for rotating refresh token mechanism ----
// func (j *JWT) RenewToken(tokenString string) (*Claims, error) {
//     claims := &Claims{}

//     token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//         return j.Secret, nil
//     }, jwt.WithoutClaimsValidation()) 

// 	 claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.RefreshTokenTTL))

//     if err != nil {
//         if errors.Is(err, jwt.ErrTokenExpired) {
// 				return  claims, nil
//         } else {
//             logger.Error("jwt.RenewToken: Failed to parse token", "error", err)
//             return nil, err
//         }
//     }

// 	if !token.Valid {
// 		logger.Error("jwt.RenewToken: Invalid token signature")
// 		return nil, fmt.Errorf("invalid token")
// 	}


//     return claims, nil
// }


func (j *JWT) VerifyToken(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            logger.Error("jwt.VerifyToken: Unexpected signing method", "method", token.Header["alg"])
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return j.Secret, nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            logger.Error("jwt.VerifyToken: Token expired", "error", err)
            return nil, errx.ErrExpiredBearerToken
        } else if errors.Is(err, jwt.ErrTokenMalformed) {
            logger.Error("jwt.VerifyToken: Token malformed", "error", err)
            return nil, errx.ErrxMalformedBearerToken
        } else {
            logger.Error("jwt.VerifyToken: JWT parse error", "error", err)
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