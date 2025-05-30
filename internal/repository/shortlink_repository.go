package repository

import (
	"javaneseivankov/url-shortener/internal/errx"
	"time"

	"github.com/google/uuid"
)

type ShortLink struct {
	 Id uuid.UUID
    ShortName  string
    OriginalUrl string
    CreatedAt  time.Time
    UserId     uuid.UUID
}

type IShortLinkRepository interface {
    CreateRedirectLink(shortName string, shortLink ShortLink) (*ShortLink, error)
    GetRedirectLink(shortName string) (*ShortLink, error)
    DeleteRedirectLink(shortName string, userId uuid.UUID) error
    EditShortLink(shortLink string, url string, userId  uuid.UUID) (*ShortLink, error)
}

type ShortLinkImpl struct {
    store map[string]ShortLink
}

func NewShortLinkRepository() IShortLinkRepository {
    return &ShortLinkImpl{
        store: make(map[string]ShortLink),
    }
}

func (s *ShortLinkImpl) CreateRedirectLink(shortName string, shortLink ShortLink) (*ShortLink, error) {
    sLink, exists := s.store[shortName]
    if exists {
        return nil, errx.ErrShortLinkAlreadyExists
    }

    s.store[shortName] = shortLink

    res := &ShortLink{
        Id:         sLink.Id,
        ShortName:  shortName,
        OriginalUrl: shortLink.OriginalUrl,
        CreatedAt:  time.Now(),
    }

    return res, nil
}

func (s *ShortLinkImpl) GetRedirectLink(shortName string) (*ShortLink, error) {
    sLink, exists := s.store[shortName]
    if !exists {
        return nil, errx.ErrShortLinkNotFound
    }

    res := &ShortLink{
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

func (s *ShortLinkImpl) EditShortLink(shortName string, url string, userId uuid.UUID) (*ShortLink, error) {
    sLink, exists := s.store[shortName]
    if !exists {
        return nil, errx.ErrShortLinkNotFound
    }

	 if userId !=  sLink.UserId {
		return nil, errx.ErrShortLinkUnauthorizedOperation
	 }


    res := &ShortLink{
        Id:         sLink.Id,
        ShortName:  shortName,
        OriginalUrl: url,
        CreatedAt:  time.Now(),
    }

    s.store[shortName] = *res

    return res, nil
}
