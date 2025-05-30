package rest

import (
	"encoding/json"
	"javaneseivankov/url-shortener/internal/app_errors"
	"javaneseivankov/url-shortener/internal/dto"
	"javaneseivankov/url-shortener/internal/service"
	"javaneseivankov/url-shortener/pkg"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: Move this to appropriate module
func ApiPathV1(path string) string {
	return "/api/v1" + path 
}


type ShortLinkController struct {
	svc service.IShortLinkService
}

func NewShortLinkController(service service.IShortLinkService) *ShortLinkController {
	return &ShortLinkController{
		svc: service,
	}
}

func (s *ShortLinkController) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestShortenLink

	ctx := r.Context()
	claims, ok := ctx.Value("claims").(pkg.Claims)
	if !ok {
		// TODO: Add logging
		pkg.SendError(w, app_errors.ErrInternalServerError)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.SendError(w, err)
		return;
	}

	res, err := s.svc.CreateShortLink(req.ShortName, req.Url, &claims)
	pkg.SendResponse(w, res, http.StatusOK, err)
}

func (s *ShortLinkController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortName := vars["shortName"]

	ctx := r.Context()
	claims, ok := ctx.Value("claims").(pkg.Claims)
	if !ok {
		// TODO: Add logging
		pkg.SendError(w, app_errors.ErrInternalServerError)
	}

	err := s.svc.DeleteShortLink(shortName, &claims)
	pkg.SendResponse(w, map[string]any{}, http.StatusOK, err)
}

func (s *ShortLinkController) RedirectHandler(w http.ResponseWriter, r *http.Request) {	vars := mux.Vars(r)
	shortName := vars["shortName"]

	res, err := s.svc.GetRedirectLink(shortName)
	if err != nil {
		pkg.SendError(w, err)
	}

	http.Redirect(w, r, res.Url, http.StatusTemporaryRedirect)
}

func (s *ShortLinkController) EditShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortName := vars["shortName"]

	var req dto.RequestEditShortLink

	ctx := r.Context()
	claims, ok := ctx.Value("claims").(pkg.Claims)
	if !ok {
		// TODO: Add logging
		pkg.SendError(w, app_errors.ErrInternalServerError)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.SendError(w, err)
		return
	}

	res, err := s.svc.EditShortLink(shortName, req.NewUrl, &claims);
	pkg.SendResponse(w, res, http.StatusOK, err)
}