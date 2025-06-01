package repository

import (
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/internal/repository/model"
	"time"

	"github.com/google/uuid"
)

type IShortLinkRepository interface {
    CreateRedirectLink(shortName string, shortLink model.ShortLink) (*model.ShortLink, error)
    GetRedirectLink(shortName string) (*model.ShortLink, error)
    DeleteRedirectLink(shortName string, userId uuid.UUID) error
    EditShortLink(shortLink string, url string, userId  uuid.UUID) (*model.ShortLink, error)
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
    sLink, exists := s.store[shortName]
    if exists {
        return nil, errx.ErrShortLinkAlreadyExists
    }

    s.store[shortName] = shortLink

    res := &model.ShortLink{
        Id:         sLink.Id,
        ShortName:  shortName,
        OriginalUrl: shortLink.OriginalUrl,
        CreatedAt:  time.Now(),
    }

    return res, nil
}

func (s *ShortLinkImpl) GetRedirectLink(shortName string) (*model.ShortLink, error) {
    sLink, exists := s.store[shortName]
    if !exists {
        return nil, errx.ErrShortLinkNotFound
    }

    res := &model.ShortLink{
        Id:         sLink.Id,
        ShortName:  shortName,
        OriginalUrl: sLink.OriginalUrl,
        CreatedAt:  sLink.CreatedAt,
    }

    return res, nil
}

func (s *ShortLinkImpl) DeleteRedirectLink(shortName string, userId uuid.UUID) error {
    sLink, exists := s.store[shortName]
    if !exists {
        return errx.ErrShortLinkNotFound
    }

	 if sLink.Id != userId {
		return errx.ErrShortLinkUnauthorizedOperation
	 }

    delete(s.store, shortName)
    return nil
}

func (s *ShortLinkImpl) EditShortLink(shortName string, url string, userId uuid.UUID) (*model.ShortLink, error) {
    sLink, exists := s.store[shortName]
    if !exists {
        return nil, errx.ErrShortLinkNotFound
    }

	 if userId !=  sLink.UserId {
		return nil, errx.ErrShortLinkUnauthorizedOperation
	 }


    res := &model.ShortLink{
        Id:         sLink.Id,
        ShortName:  shortName,
        OriginalUrl: url,
        CreatedAt:  time.Now(),
    }

    s.store[shortName] = *res

    return res, nil
}
