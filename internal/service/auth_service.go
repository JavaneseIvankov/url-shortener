package service

import (
	"context"
	"errors"
	"javaneseivankov/url-shortener/internal/dto"
	"javaneseivankov/url-shortener/internal/repository"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg/bcrypt"
	"javaneseivankov/url-shortener/pkg/jwt"
	"javaneseivankov/url-shortener/pkg/logger"

	"github.com/google/uuid"
)

type AuthService struct {
    r repository.IUserRepository
    j jwt.JWT
}

func NewAuthService(userRepo repository.IUserRepository, jwtAuth jwt.JWT) *AuthService {
    return &AuthService{r: userRepo, j: jwtAuth}
}

func (auth *AuthService) RegisterUser(ctx context.Context, email string, password string) (*dto.RegisterResponse, error) {
    logger.Info("AuthService.RegisterUser: registering user", "email", email)

    hashedPassword, err := bcrypt.Hash(password)
    if err != nil {
        logger.Error("AuthService.RegisterUser: failed to hash password", "error", err)
        return nil, errors.New("failed to hash password")
    }

    userId := uuid.New()
    user := model.User{
        ID:       userId,
        Email:    email,
        Password: hashedPassword,
    }

    logger.Info("AuthService.RegisterUser: creating user in repository", "email", email)
    if err = auth.r.CreateUser(ctx, user); err != nil {
        logger.Error("AuthService.RegisterUser: failed to create user in repository", "email", email, "error", err)
        return nil, err
    }

    logger.Info("AuthService.RegisterUser: generating token for user", "email", email)
    accessToken, err := auth.j.GenerateAccessToken(&user)
    if err != nil {
        logger.Error("AuthService.RegisterUser: failed to generate access token", "email", email, "error", err)
        return nil, err
    }

    refreshToken, err := auth.j.GenerateRefreshToken(&user)
    if err != nil {
        logger.Error("AuthService.RegisterUser: failed to generate refresh token", "email", email, "error", err)
        return nil, err
    }
    res := &dto.RegisterResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }

    logger.Info("AuthService.RegisterUser: user registered successfully", "email", email)
    logger.Info("AuthService.RefreshSession: session refreshed successfully")
    logger.Debug("AuthService.RefreshSession: session refreshed successfully", "accessToken", accessToken, "refreshToken", refreshToken)
    return res, nil
}

func (auth *AuthService) LoginUser(ctx context.Context,email string, password string) (*dto.LoginResponse, error) {
    logger.Info("AuthService.LoginUser: logging in user", "email", email)

    user, err := auth.r.GetUserByEmail(ctx, email)
    if err != nil || user == nil {
        logger.Error("AuthService.LoginUser: Error getting user from repository", "email", email, "error", err)
        return nil, errors.New("invalid email or password")
    }

    if ok := bcrypt.Compare(password, user.Password); !ok {
        logger.Error("AuthService.LoginUser: invalid email or password", "email", email, "storedPassword", user.Password, "password", password)
        return nil, errors.New("invalid email or password")
    }

    logger.Info("AuthService.LoginUser: generating token for user", "email", email)
    accessToken, err := auth.j.GenerateAccessToken(user)
    if err != nil {
        logger.Error("AuthService.LoginUser: failed to generate token", "email", email, "error", err)
        return nil, err
    }

    refreshToken, err := auth.j.GenerateRefreshToken(user)
    if err != nil {
        logger.Error("AuthService.LoginUser: failed to generate refresh token", "email", email, "error", err)
        return nil, err
    }

    res := &dto.LoginResponse{
        AccessToken: accessToken,
		  RefreshToken: refreshToken,
    }

    logger.Info("AuthService.LoginUser: user logged in successfully", "email", email)
    logger.Debug("AuthService.LoginUser: user logged in successfully", "accessToken", accessToken)
    return res, nil
}

func (auth *AuthService) RefreshSession(ctx context.Context, refreshToken string) (*dto.RefreshSessionResponse, error) {
    logger.Info("AuthService.RefreshSession: refreshing session", "refreshToken", refreshToken)
	 

    accessToken, err := auth.j.RenewAccessToken(refreshToken)
    if err != nil {
        logger.Error("AuthService.RefreshSession: failed to renew access token", "refreshToken", refreshToken, "error", err)
        return nil, err
    }

    res := &dto.RefreshSessionResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }

    logger.Info("AuthService.RefreshSession: session refreshed successfully")
    logger.Debug("AuthService.RefreshSession: session refreshed successfully", "accessToken", accessToken, "refreshToken", refreshToken)
    return res, nil
}