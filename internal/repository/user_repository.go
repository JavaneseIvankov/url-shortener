package repository

import (
	"errors"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg"
	"sync"

	"github.com/google/uuid"
)

type IUserRepository interface {
    CreateUser(user model.User ) error
    UpdateUser(user model.User) error
    GetUserByEmail(email string) (model.User, error)
    GetUserByID(id uuid.UUID) (model.User, error)
}

type UserRepositoryImpl struct {
    usersByID    map[uuid.UUID]model.User
    usersByEmail map[string]model.User
    mu           sync.Mutex
}

func NewUserRepository() IUserRepository {
    return &UserRepositoryImpl{
        usersByID:    make(map[uuid.UUID]model.User),
        usersByEmail: make(map[string]model.User),
    }
}

func (r *UserRepositoryImpl) CreateUser(user model.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.usersByID[user.ID]; exists {
        pkg.Error("CreateUser: user with this ID already exists", "userID", user.ID)
        return errors.New("user with this ID already exists")
    }
    if _, exists := r.usersByEmail[user.Email]; exists {
        pkg.Error("CreateUser: user with this email already exists", "email", user.Email)
        return errors.New("user with this email already exists")
    }

    r.usersByID[user.ID] = user
    r.usersByEmail[user.Email] = user
    pkg.Info("CreateUser: user created successfully", "userID", user.ID, "email", user.Email)
    return nil
}

func (r *UserRepositoryImpl) UpdateUser(user model.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.usersByID[user.ID]; !exists {
        pkg.Error("UpdateUser: user with this ID does not exist", "userID", user.ID)
        return errors.New("user with this ID does not exist")
    }

    r.usersByID[user.ID] = user
    r.usersByEmail[user.Email] = user
    pkg.Info("UpdateUser: user updated successfully", "userID", user.ID, "email", user.Email)
    return nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (model.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    user, exists := r.usersByEmail[email]
    if !exists {
        pkg.Error("GetUserByEmail: user with this email does not exist", "email", email)
        return model.User{}, errors.New("user with this email does not exist")
    }
    pkg.Info("GetUserByEmail: user retrieved successfully", "email", email)
    return user, nil
}

func (r *UserRepositoryImpl) GetUserByID(id uuid.UUID) (model.User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    user, exists := r.usersByID[id]
    if !exists {
        pkg.Error("GetUserByID: user with this ID does not exist", "userID", id)
        return model.User{}, errors.New("user with this ID does not exist")
    }
    pkg.Info("GetUserByID: user retrieved successfully", "userID", id)
    return user, nil
}