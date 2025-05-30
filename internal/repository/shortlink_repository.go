package repository

import (
	"fmt"
	"javaneseivankov/url-shortener/internal/app_errors"
	"time"
)

func uuid() uint {
	return 10
}

type ShortLink struct {
	id uint
	ShortName string
	OriginalUrl string
	createdAt string
}

type IShortLinkRepository interface {
	CreateRedirectLink(shortName string, url string) (*ShortLink, error)
	GetRedirectLink(shortName string) (*ShortLink, error)
	DeleteRedirectLink(shortName string) error
	EditShortLink(shortName string, newUrl string) (*ShortLink, error)
}

type ShortLinkImpl struct {
	store map[string]string 
}

func NewShortLinkRepository() IShortLinkRepository {
	return &ShortLinkImpl{
		store: make(map[string]string),
}
} 

func (s *ShortLinkImpl) CreateRedirectLink(shortName string, url string) (*ShortLink, error) {
	_, exists := s.store[shortName]
	if exists {
		fmt.Println("LINK EXISTS")
		return nil, app_errors.ErrShortLinkAlreadyExists
	}

	s.store[shortName] = url
	
	res := &ShortLink{
		id: uuid(),
		ShortName: shortName,
		OriginalUrl:  url,
		createdAt: time.Now().Format("2006-01-02"),
	}

	return res, nil
}


func (s *ShortLinkImpl) GetRedirectLink(shortName string) (*ShortLink, error) {
	url, exists := s.store[shortName]
	if !exists {
		return nil, app_errors.ErrShortLinkAlreadyExists
	}
	
	res := &ShortLink{
		id: uuid(),
		ShortName: shortName,
		OriginalUrl:  url,
		createdAt: time.Now().Format("2006-12-06"),
	}

	return res, nil
}

func (s *ShortLinkImpl) DeleteRedirectLink(shortName string) (error) {
	_, exists := s.store[shortName]
	if !exists {
		return app_errors.ErrShortLinkNotFound
	}

	delete(s.store, shortName)
	return nil;
}

func (s *ShortLinkImpl) EditShortLink(shortName string, newUrl string) (*ShortLink, error) {
	_, exists := s.store[shortName]
	if !exists {
		return nil, app_errors.ErrShortLinkNotFound
	}

	s.store[shortName] = newUrl
	res := &ShortLink{
		id: uuid(),
		ShortName: shortName,
		OriginalUrl: newUrl,
		createdAt: time.Now().Format("2006-12-06"),
	}

	return res, nil
}
