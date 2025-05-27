package repostiories

import (
	"errors"
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

type IShortlinkRepository interface {
	CreateRedirectLink(shortName string, url string) (*ShortLink, error)
	GetRedirectLink(shortName string) (*ShortLink, error)
	DeleteRedirectLink(shortName string) error
}

type ShortLinkImpl struct {
	store map[string]string 
}

func NewShortLinkRepository() IShortlinkRepository {
	return &ShortLinkImpl{
		store: make(map[string]string),
	}
} 

func (s *ShortLinkImpl) CreateRedirectLink(shortName string, url string) (*ShortLink, error) {
	_, exists := s.store[shortName]
	if exists {
		return nil, errors.New("Link already exists")
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
		return nil, errors.New("Link doesn't exist")
	}
	
	res := &ShortLink{
		id: uuid(),
		ShortName: shortName,
		OriginalUrl:  url,
		createdAt: time.Now().Format("YYYY-MM-DD"),
	}

	return res, nil
}

func (s *ShortLinkImpl) DeleteRedirectLink(shortName string) (error) {
	_, exists := s.store[shortName]
	if !exists {
		return errors.New("Link with name " + shortName + "doesn't exist")
	}

	delete(s.store, shortName)
	return nil;
}
