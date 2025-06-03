package repository

import (
	"context"
	"javaneseivankov/url-shortener/internal/repository/model"

	"github.com/google/uuid"
)

type IShortLinkRepository interface {
    CreateRedirectLink(ctx context.Context, shortLink model.ShortLink) (*model.ShortLink, error)
    GetRedirectLink(ctx context.Context, shortName string) (*model.ShortLink, error)
    DeleteRedirectLink(ctx context.Context, shortName string, userId uuid.UUID) error
    EditShortLink(ctx context.Context, shortName string, url string, userId uuid.UUID) (*model.ShortLink, error)
}
