package rest

import (
	"encoding/json"
	"net/http"

	"javaneseivankov/url-shortener/internal/dto"
	"javaneseivankov/url-shortener/internal/service"
	"javaneseivankov/url-shortener/pkg"
	"javaneseivankov/url-shortener/pkg/logger"
)

type AuthController struct {
    authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
    return &AuthController{
        authService: authService,
    }
}

func (c *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
    logger.Info("AuthController.RegisterUser: handling user registration")

    var req dto.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("AuthController.RegisterUser: failed to decode request body", "error", err)
        pkg.SendError(w, err)
        return
    }

    logger.Info("AuthController.RegisterUser: calling authService.RegisterUser", "email", req.Email)
    res, err := c.authService.RegisterUser(r.Context() ,req.Email, req.Password)
    if err != nil {
        logger.Error("AuthController.RegisterUser: failed to register user", "email", req.Email, "error", err)
        pkg.SendError(w, err)
        return
    }

    logger.Info("AuthController.RegisterUser: user registered successfully", "email", req.Email)
    pkg.SendResponse(w, res, http.StatusOK, err)
}

func (c *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
    logger.Info("AuthController.LoginUser: handling user login")

    var req dto.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("AuthController.LoginUser: failed to decode request body", "error", err)
        pkg.SendError(w, err)
        return
    }

    logger.Info("AuthController.LoginUser: calling authService.LoginUser", "email", req.Email)
    res, err := c.authService.LoginUser(r.Context(), req.Email, req.Password)
    if err != nil {
        logger.Error("AuthController.LoginUser: failed to login user", "email", req.Email, "error", err)
        pkg.SendError(w, err)
        return
    }

    logger.Info("AuthController.LoginUser: user logged in successfully", "email", req.Email)
    pkg.SendResponse(w, res, http.StatusOK, err)
}

func (c *AuthController) RefreshSession(w http.ResponseWriter, r *http.Request) {
	logger.Info("AuthController.RefreshSession: refreshing session...")

	var req dto.RefreshSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.Error("AuthController.RefreshSession: failed to decode request body", "error", err)
        pkg.SendError(w, err)
        return
	}

	res, err := c.authService.RefreshSession(r.Context(), req.RefreshToken)
	if err != nil {
		logger.Error("AuthController.RefreshSession: failed to refresh session", "error", err)
		pkg.SendError(w, err)
	}

	pkg.SendResponse(w, res, http.StatusOK, err)
}