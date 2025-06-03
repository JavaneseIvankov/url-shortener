package rest

import (
	"encoding/json"
	"javaneseivankov/url-shortener/internal/contextkeys"
	"javaneseivankov/url-shortener/internal/dto"
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/internal/service"
	"javaneseivankov/url-shortener/pkg"
	"javaneseivankov/url-shortener/pkg/jwt"
	"javaneseivankov/url-shortener/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type ShortLinkController struct {
    svc service.IShortLinkService
}

func NewShortLinkController(service service.IShortLinkService) *ShortLinkController {
    return &ShortLinkController{
        svc: service,
    }
}

func (s *ShortLinkController) ShortenHandler(w http.ResponseWriter, r *http.Request) {
    logger.Info("ShortLinkController.ShortenHandler: handling short link creation")

    var req dto.RequestShortenLink

    ctx := r.Context()
    claims, ok := ctx.Value(contextkeys.ClaimCtxKey).(*jwt.Claims)
    if !ok {
        logger.Error("ShortLinkController.ShortenHandler: failed to parse claims from context")
        pkg.SendError(w, errx.ErrInternalServerError)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("ShortLinkController.ShortenHandler: failed to decode request body", "error", err)
        pkg.SendError(w, err)
        return
    }

    logger.Info("ShortLinkController.ShortenHandler: calling service.CreateShortLink", "shortName", req.ShortName, "url", req.Url)
    res, err := s.svc.CreateShortLink(r.Context(), req.ShortName, req.Url, claims)
    if err != nil {
        logger.Error("ShortLinkController.ShortenHandler: failed to create short link", "shortName", req.ShortName, "error", err)
    }
    pkg.SendResponse(w, res, http.StatusOK, err)
}

func (s *ShortLinkController) DeleteShortLinkHandler(w http.ResponseWriter, r *http.Request) {
    logger.Info("ShortLinkController.DeleteShortLinkHandler: handling short link deletion")

    vars := mux.Vars(r)
    shortName := vars["shortName"]

    ctx := r.Context()
    claims, ok := ctx.Value(contextkeys.ClaimCtxKey).(*jwt.Claims)
    if !ok {
        logger.Error("ShortLinkController.DeleteShortLinkHandler: failed to parse claims from context")
        pkg.SendError(w, errx.ErrInternalServerError)
        return
    }

    logger.Info("ShortLinkController.DeleteShortLinkHandler: calling service.DeleteShortLink", "shortName", shortName)
    err := s.svc.DeleteShortLink(r.Context(), shortName, claims)
    if err != nil {
        logger.Error("ShortLinkController.DeleteShortLinkHandler: failed to delete short link", "shortName", shortName, "error", err)
    }
    pkg.SendResponse(w, map[string]any{}, http.StatusOK, err)
}

func (s *ShortLinkController) RedirectHandler(w http.ResponseWriter, r *http.Request) {
    logger.Info("ShortLinkController.RedirectHandler: handling redirect")

    vars := mux.Vars(r)
    shortName := vars["shortName"]

    logger.Info("ShortLinkController.RedirectHandler: calling service.GetRedirectLink", "shortName", shortName)
    res, err := s.svc.GetRedirectLink(r.Context(), shortName)
    if err != nil {
        logger.Error("ShortLinkController.RedirectHandler: failed to retrieve redirect link", "shortName", shortName, "error", err)
        pkg.SendError(w, err)
        return
    }

    logger.Info("ShortLinkController.RedirectHandler: redirecting to URL", "url", res.Url)
    http.Redirect(w, r, res.Url, http.StatusTemporaryRedirect)
}

func (s *ShortLinkController) EditShortLinkHandler(w http.ResponseWriter, r *http.Request) {
    logger.Info("ShortLinkController.EditShortLinkHandler: handling short link editing")

    vars := mux.Vars(r)
    shortName := vars["shortName"]

    var req dto.RequestEditShortLink

    ctx := r.Context()
    claims, ok := ctx.Value(contextkeys.ClaimCtxKey).(*jwt.Claims)
    if !ok {
        logger.Error("ShortLinkController.EditShortLinkHandler: failed to parse claims from context")
        pkg.SendError(w, errx.ErrInternalServerError)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("ShortLinkController.EditShortLinkHandler: failed to decode request body", "error", err)
        pkg.SendError(w, err)
        return
    }

    logger.Info("ShortLinkController.EditShortLinkHandler: calling service.EditShortLink", "shortName", shortName, "newUrl", req.NewUrl)
    res, err := s.svc.EditShortLink(r.Context(), shortName, req.NewUrl, claims)
    if err != nil {
        logger.Error("ShortLinkController.EditShortLinkHandler: failed to edit short link", "shortName", shortName, "error", err)
    }
    pkg.SendResponse(w, res, http.StatusOK, err)
}