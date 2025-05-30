package service

import (
	"javaneseivankov/url-shortener/internal/dto"
	"javaneseivankov/url-shortener/internal/repository"
)

type IShortLinkService interface {
	CreateShortLink(shortName string, url string) (*dto.ResponseShortenLink, error)
	EditShortLink(shortName string, url string) (*dto.ResponseEditShortLink, error)
	GetRedirectLink(shortName string) (*dto.ResponseGetShortLink, error)
	DeleteShortLink(shortName string) error 
}

type ShortLinkService struct {
	repo repository.IShortLinkRepository
}

func NewShortLinkService(repository repository.IShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{
		repo: repository,
	}
}

func (s *ShortLinkService) CreateShortLink(shortName string, url string) (*dto.ResponseShortenLink, error) {
	repoRes, err := s.repo.CreateRedirectLink(shortName, url);

	 if err != nil {
		return nil, err
	} 
	redirectUrl :=  "/s/" + repoRes.ShortName

	res := dto.ResponseShortenLink{
		Url: redirectUrl,
	}
	return &res, nil
}

func (s *ShortLinkService) EditShortLink(shortName string, url string) (*dto.ResponseEditShortLink, error) {
	repoRes, err := s.repo.EditShortLink(shortName, url)

	if err != nil {
		return nil, err
	}

	res := dto.ResponseEditShortLink{
		Url: repoRes.OriginalUrl,
		ShortName: repoRes.ShortName,
	}

	return &res, nil
}


func (s *ShortLinkService) GetRedirectLink(shortName string) (*dto.ResponseGetShortLink, error) {
	slink, err := s.repo.GetRedirectLink(shortName)

	if err != nil {
		return nil, err
	}

	res := dto.ResponseGetShortLink{
		Url: slink.OriginalUrl,
	}

	return &res, nil
}

func (s *ShortLinkService) DeleteShortLink(shortName string) error {
	err := s.repo.DeleteRedirectLink(shortName)
	return err;
}