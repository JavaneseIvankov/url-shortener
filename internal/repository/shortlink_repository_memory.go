package repository

import (
	"context"
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg/logger"
	"time"

	"github.com/google/uuid"
)


type ShortLinkRepositoryMemory struct {
    store map[string]model.ShortLink
}

func NewShortLinkRepositoryMemory() IShortLinkRepository {
    return &ShortLinkRepositoryMemory{
        store: make(map[string]model.ShortLink),
    }
}

func (s *ShortLinkRepositoryMemory) CreateRedirectLink(ctx context.Context, shortLink model.ShortLink) (*model.ShortLink, error) {
    _, exists := s.store[shortLink.ShortName]
    if exists {
        logger.Error("ShortLinkRepository.CreateRedirectLink: short link already exists", "shortName", shortLink.ShortName)
        return nil, errx.ErrShortLinkAlreadyExists
    }

    s.store[shortLink.ShortName] = shortLink
    logger.Info("ShortLinkRepository.CreateRedirectLink: short link created successfully", "shortName", shortLink.ShortName, "originalUrl", shortLink.OriginalUrl)

    res := &model.ShortLink{
        Id:         shortLink.Id,
        ShortName:  shortLink.ShortName,
        OriginalUrl: shortLink.OriginalUrl,
        CreatedAt:  time.Now(),
    }

    return res, nil
}

func (s *ShortLinkRepositoryMemory) GetRedirectLink(ctx context.Context, shortName string) (*model.ShortLink, error) {
    sLink, exists := s.store[shortName]
    if !exists {
        logger.Error("ShortLinkRepository.GetRedirectLink: short link not found", "shortName", shortName)
        return nil, errx.ErrShortLinkNotFound
    }

    logger.Info("ShortLinkRepository.GetRedirectLink: short link retrieved successfully", "shortName", shortName, "originalUrl", sLink.OriginalUrl)
    return &sLink, nil
}

func (s *ShortLinkRepositoryMemory) DeleteRedirectLink(ctx context.Context, shortName string, userId uuid.UUID) error {
    sLink, exists := s.store[shortName]
    if !exists {
        logger.Error("ShortLinkRepository.DeleteRedirectLink: short link not found", "shortName", shortName)
        return errx.ErrShortLinkNotFound
    }

    if sLink.UserId != userId {
        logger.Error("ShortLinkRepository.DeleteRedirectLink: unauthorized operation", "shortName", shortName, "userId", userId)
        return errx.ErrShortLinkUnauthorizedOperation
    }

    delete(s.store, shortName)
    logger.Info("ShortLinkRepository.DeleteRedirectLink: short link deleted successfully", "shortName", shortName, "userId", userId)
    return nil
}

func (s *ShortLinkRepositoryMemory) EditShortLink(ctx context.Context, shortName string, url string, userId uuid.UUID) (*model.ShortLink, error) {
    sLink, exists := s.store[shortName]
    if !exists {
        logger.Error("ShortLinkRepository.EditShortLink: short link not found", "shortName", shortName)
        return nil, errx.ErrShortLinkNotFound
    }

    if userId != sLink.UserId {
        logger.Error("ShortLinkRepository.EditShortLink: unauthorized operation", "shortName", shortName, "userId", userId)
        return nil, errx.ErrShortLinkUnauthorizedOperation
    }

    res := &model.ShortLink{
        Id:         sLink.Id,
        ShortName:  shortName,
        OriginalUrl: url,
        CreatedAt:  time.Now(),
    }

    s.store[shortName] = *res
    logger.Info("ShortLinkRepository.EditShortLink: short link updated successfully", "shortName", shortName, "originalUrl", url, "userId", userId)
    return res, nil
}