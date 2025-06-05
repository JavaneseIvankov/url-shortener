package repository

import (
	"context"
	"javaneseivankov/url-shortener/internal/repository/model"

	"github.com/google/uuid"
)

type IUserRepository interface {
    CreateUser(ctx context.Context, user model.User ) error
    UpdateUser(ctx context.Context,user model.User) error
    GetUserByEmail(ctx context.Context, email string) (*model.User, error)
    GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}
