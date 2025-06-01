package service

import (
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

func (auth *AuthService) RegisterUser(email string, password string) (*dto.RegisterResponse, error) {
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
    if err = auth.r.CreateUser(user); err != nil {
        logger.Error("AuthService.RegisterUser: failed to create user in repository", "email", email, "error", err)
        return nil, err
    }

    logger.Info("AuthService.RegisterUser: generating token for user", "email", email)
    token, err := auth.j.GenerateToken(user)
    if err != nil {
        logger.Error("AuthService.RegisterUser: failed to generate token", "email", email, "error", err)
        return nil, err
    }

    res := &dto.RegisterResponse{
        AccessToken:  token,
        RefreshToken: token,
    }

    logger.Info("AuthService.RegisterUser: user registered successfully", "email", email)
    return res, nil
}

func (auth *AuthService) LoginUser(email string, password string) (*dto.LoginResponse, error) {
    logger.Info("AuthService.LoginUser: logging in user", "email", email)

    user, err := auth.r.GetUserByEmail(email)
    if err != nil {
        logger.Error("AuthService.LoginUser: invalid email or password", "email", email, "error", err)
        return nil, errors.New("invalid email or password")
    }

    if ok := bcrypt.Compare(user.Password, password); !ok {
        logger.Error("AuthService.LoginUser: invalid email or password", "email", email)
        return nil, errors.New("invalid email or password")
    }

    logger.Info("AuthService.LoginUser: generating token for user", "email", email)
    token, err := auth.j.GenerateToken(user)
    if err != nil {
        logger.Error("AuthService.LoginUser: failed to generate token", "email", email, "error", err)
        return nil, err
    }

    res := &dto.LoginResponse{
        AccessToken: token,
    }

    logger.Info("AuthService.LoginUser: user logged in successfully", "email", email)
    return res, nil
}