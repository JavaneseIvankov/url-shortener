package repository

import (
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type IShortLinkRepository interface {
    CreateRedirectLink(shortName string, shortLink model.ShortLink) (*model.ShortLink, error)
    GetRedirectLink(shortName string) (*model.ShortLink, error)
    DeleteRedirectLink(shortName string, userId uuid.UUID) error
    EditShortLink(shortName string, url string, userId uuid.UUID) (*model.ShortLink, error)
}

type ShortLinkImpl struct {
    store map[string]model.ShortLink
}

func NewShortLinkRepository() IShortLinkRepository {
    return &ShortLinkImpl{
        store: make(map[string]model.ShortLink),
    }
}

func (s *ShortLinkImpl) CreateRedirectLink(shortName string, shortLink model.ShortLink) (*model.ShortLink, error) {
    _, exists := s.store[shortName]
    if exists {
        logger.Error("ShortLinkRepository.CreateRedirectLink: short link already exists", "shortName", shortName)
        return nil, errx.ErrShortLinkAlreadyExists
    }

    s.store[shortName] = shortLink
    logger.Info("ShortLinkRepository.CreateRedirectLink: short link created successfully", "shortName", shortName, "originalUrl", shortLink.OriginalUrl)

    res := &model.ShortLink{
        Id:         shortLink.Id,
        ShortName:  shortName,
        OriginalUrl: shortLink.OriginalUrl,
        CreatedAt:  time.Now(),
    }

    return res, nil
}

func (s *ShortLinkImpl) GetRedirectLink(shortName string) (*model.ShortLink, error) {
    sLink, exists := s.store[shortName]
    if !exists {
        logger.Error("ShortLinkRepository.GetRedirectLink: short link not found", "shortName", shortName)
        return nil, errx.ErrShortLinkNotFound
    }

    logger.Info("ShortLinkRepository.GetRedirectLink: short link retrieved successfully", "shortName", shortName, "originalUrl", sLink.OriginalUrl)
    return &sLink, nil
}

func (s *ShortLinkImpl) DeleteRedirectLink(shortName string, userId uuid.UUID) error {
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

func (s *ShortLinkImpl) EditShortLink(shortName string, url string, userId uuid.UUID) (*model.ShortLink, error) {
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