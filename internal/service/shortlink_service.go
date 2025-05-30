package service

import (
	"javaneseivankov/url-shortener/internal/dto"
	"javaneseivankov/url-shortener/internal/repository"
	"javaneseivankov/url-shortener/pkg"

	"github.com/google/uuid"
)

type IShortLinkService interface {
	CreateShortLink(shortName string, url string, claims *pkg.Claims) (*dto.ResponseShortenLink, error)
	EditShortLink(shortName string, url string, claims *pkg.Claims) (*dto.ResponseEditShortLink, error)
	GetRedirectLink(shortName string) (*dto.ResponseGetShortLink, error)
	DeleteShortLink(shortName string, claims *pkg.Claims) error 
}

type ShortLinkService struct {
	repo repository.IShortLinkRepository
}

func NewShortLinkService(repository repository.IShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{
		repo: repository,
	}
}

func (s *ShortLinkService) CreateShortLink(shortName string, url string, claims *pkg.Claims) (*dto.ResponseShortenLink, error) {
   generatedId, err := uuid.NewUUID()
   if err != nil {
      return nil, err
   }

   sLink := repository.ShortLink{
      Id: generatedId,
      ShortName: shortName,
      OriginalUrl: url,
      UserId: claims.UserID,
   }

	repoRes, err := s.repo.CreateRedirectLink(shortName, sLink);

	 if err != nil {
		return nil, err
	} 
	redirectUrl :=  "/s/" + repoRes.ShortName

	res := dto.ResponseShortenLink{
		Url: redirectUrl,
	}
	return &res, nil
}

func (s *ShortLinkService) EditShortLink(shortName string, url string, claims *pkg.Claims) (*dto.ResponseEditShortLink, error) {
   sLink := repository.ShortLink{
      ShortName: shortName,
      OriginalUrl: url,
   }
   
	_, err := s.repo.EditShortLink(shortName, url, claims.UserID)

	if err != nil {
		return nil, err
	}

	res := dto.ResponseEditShortLink{
		Url: sLink.OriginalUrl,
		ShortName: sLink.ShortName,
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

func (s *ShortLinkService) DeleteShortLink(shortName string, claims *pkg.Claims) error {
	err := s.repo.DeleteRedirectLink(shortName, claims.UserID)
	return err;
}