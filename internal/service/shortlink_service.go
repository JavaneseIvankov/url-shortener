package service

import (
	"context"
	"javaneseivankov/url-shortener/internal/dto"
	"javaneseivankov/url-shortener/internal/repository"
	"javaneseivankov/url-shortener/internal/repository/model"
	"javaneseivankov/url-shortener/pkg/jwt"
	"javaneseivankov/url-shortener/pkg/logger"

	"github.com/google/uuid"
)

type IShortLinkService interface {
    CreateShortLink(ctx context.Context, shortName string, url string, claims *jwt.Claims) (*dto.ResponseShortenLink, error)
    EditShortLink(ctx context.Context, shortName string, url string, claims *jwt.Claims) (*dto.ResponseEditShortLink, error)
    GetRedirectLink(ctx context.Context, shortName string) (*dto.ResponseGetShortLink, error)
    DeleteShortLink(ctx context.Context, shortName string, claims *jwt.Claims) error
}

type ShortLinkService struct {
    repo repository.IShortLinkRepository
}

func NewShortLinkService(repository repository.IShortLinkRepository) IShortLinkService {
    return &ShortLinkService{
        repo: repository,
    }
}

func (s *ShortLinkService) CreateShortLink(ctx context.Context, shortName string, url string, claims *jwt.Claims) (*dto.ResponseShortenLink, error) {
    logger.Info("ShortLinkService.CreateShortLink: creating short link", "shortName", shortName, "url", url)

    generatedId, err := uuid.NewUUID()
    if err != nil {
        logger.Error("ShortLinkService.CreateShortLink: failed to generate UUID", "error", err)
        return nil, err
    }

    userId, err := uuid.Parse(claims.UserID)
    if err != nil {
        logger.Error("ShortLinkService.CreateShortLink: failed to parse user ID", "userID", claims.UserID, "error", err)
        return nil, err
    }

    sLink := model.ShortLink{
        Id:         generatedId,
        ShortName:  shortName,
        OriginalUrl: url,
        UserId:     userId,
    }

    repoRes, err := s.repo.CreateRedirectLink(ctx, sLink)
    if err != nil {
        logger.Error("ShortLinkService.CreateShortLink: failed to create redirect link in repository", "shortName", shortName, "error", err)
        return nil, err
    }

    redirectUrl := "/s/" + repoRes.ShortName
    res := dto.ResponseShortenLink{
        Url: redirectUrl,
    }

    logger.Info("ShortLinkService.CreateShortLink: short link created successfully", "shortName", shortName, "redirectUrl", redirectUrl)
    return &res, nil
}

func (s *ShortLinkService) EditShortLink(ctx context.Context, shortName string, url string, claims *jwt.Claims) (*dto.ResponseEditShortLink, error) {
    logger.Info("ShortLinkService.EditShortLink: editing short link", "shortName", shortName, "url", url)

    userId, err := uuid.Parse(claims.UserID)
    if err != nil {
        logger.Error("ShortLinkService.EditShortLink: failed to parse user ID", "userID", claims.UserID, "error", err)
        return nil, err
    }

    _, err = s.repo.EditShortLink(ctx, shortName, url, userId)
    if err != nil {
        logger.Error("ShortLinkService.EditShortLink: failed to edit short link in repository", "shortName", shortName, "error", err)
        return nil, err
    }

    res := dto.ResponseEditShortLink{
        Url:       url,
        ShortName: shortName,
    }

    logger.Info("ShortLinkService.EditShortLink: short link edited successfully", "shortName", shortName, "url", url)
    return &res, nil
}

func (s *ShortLinkService) GetRedirectLink(ctx context.Context, shortName string) (*dto.ResponseGetShortLink, error) {
    logger.Info("ShortLinkService.GetRedirectLink: retrieving redirect link", "shortName", shortName)

    slink, err := s.repo.GetRedirectLink(ctx, shortName)
    if err != nil {
        logger.Error("ShortLinkService.GetRedirectLink: failed to retrieve redirect link from repository", "shortName", shortName, "error", err)
        return nil, err
    }

    res := dto.ResponseGetShortLink{
        Url: slink.OriginalUrl,
    }

    logger.Info("ShortLinkService.GetRedirectLink: redirect link retrieved successfully", "shortName", shortName, "url", slink.OriginalUrl)
    return &res, nil
}

func (s *ShortLinkService) DeleteShortLink(ctx context.Context, shortName string, claims *jwt.Claims) error {
    logger.Info("ShortLinkService.DeleteShortLink: deleting short link", "shortName", shortName)

    userId, err := uuid.Parse(claims.UserID)
    if err != nil {
        logger.Error("ShortLinkService.DeleteShortLink: failed to parse user ID", "userID", claims.UserID, "error", err)
        return err
    }

    err = s.repo.DeleteRedirectLink(ctx, shortName, userId)
    if err != nil {
        logger.Error("ShortLinkService.DeleteShortLink: failed to delete short link in repository", "shortName", shortName, "error", err)
        return err
    }

    logger.Info("ShortLinkService.DeleteShortLink: short link deleted successfully", "shortName", shortName)
    return nil
}