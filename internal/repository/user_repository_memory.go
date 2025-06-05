package repository

import (
	"context"
	"errors"
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg/logger"
	"sync"

	"github.com/google/uuid"
)

type UserRepositoryMemory struct {
    usersByID    map[uuid.UUID]model.User
    usersByEmail map[string]model.User
    mu           sync.Mutex
}

func NewUserRepositoryMemory() IUserRepository {
    return &UserRepositoryMemory{
        usersByID:    make(map[uuid.UUID]model.User),
        usersByEmail: make(map[string]model.User),
    }
}

func (r *UserRepositoryMemory) CreateUser(ctx context.Context ,user model.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.usersByID[user.ID]; exists {
        logger.Error("UserRepository.CreateUser: user with this ID already exists", "userID", user.ID)
        return errors.New("user with this ID already exists")
    }
    if _, exists := r.usersByEmail[user.Email]; exists {
        logger.Error("UserRepository.CreateUser: user with this email already exists", "email", user.Email)
        return errors.New("user with this email already exists")
    }

    r.usersByID[user.ID] = user
    r.usersByEmail[user.Email] = user
    logger.Info("UserRepository.CreateUser: user created successfully", "userID", user.ID, "email", user.Email)
    return nil
}

func (r *UserRepositoryMemory) UpdateUser(ctx context.Context, user model.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.usersByID[user.ID]; !exists {
        logger.Error("UserRepository.UpdateUser: user with this ID does not exist", "userID", user.ID)
        return errors.New("user with this ID does not exist")
    }

    r.usersByID[user.ID] = user
    r.usersByEmail[user.Email] = user
    logger.Info("UserRepository.UpdateUser: user updated successfully", "userID", user.ID, "email", user.Email)
    return nil
}

func (r *UserRepositoryMemory) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    user, exists := r.usersByEmail[email]
    if !exists {
        logger.Error("UserRepository.GetUserByEmail: user with this email does not exist", "email", email)
        return nil, errx.ErrEmailDoesntExist
    }
    logger.Info("UserRepository.GetUserByEmail: user retrieved successfully", "email", email)
    return &user, nil
}

func (r *UserRepositoryMemory) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    user, exists := r.usersByID[id]
    if !exists {
        logger.Error("UserRepository.GetUserByID: user with this ID does not exist", "userID", id)
        return nil, errx.ErrUserIdDoesntExist
    }
    logger.Info("UserRepository.GetUserByID: user retrieved successfully", "userID", id)
    return &user, nil
}