package main

import (
	"javaneseivankov/url-shortener/internal/controller/rest"
	"javaneseivankov/url-shortener/internal/middleware"
	"javaneseivankov/url-shortener/internal/repository"
	"javaneseivankov/url-shortener/internal/service"
	"javaneseivankov/url-shortener/pkg"
	"javaneseivankov/url-shortener/pkg/jwt"
	"javaneseivankov/url-shortener/pkg/logger"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	pkg.SendJSON(w, http.StatusOK, map[string]interface{}{"payload": "pong"})
}

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtTTL := os.Getenv("JWT_TTL")

	if jwtSecret == "" || jwtTTL == "" {
		log.Fatalln("JWT_SECRET and JWT_TTL must be set")
	}

	jwtAuth := jwt.NewJWT(jwtSecret, jwtTTL)
	// TODO: Get from .env
	logger.Init("development")
	
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	requireAuth := middleware.AuthMiddleware(jwtAuth)
	apply := middleware.ApplyMiddleware

	apiRouter := router.PathPrefix("/api").Subrouter()
	v1 := apiRouter.PathPrefix("/v1").Subrouter()

	shortLinkRepo := repository.NewShortLinkRepository()
	shortLinkService := service.NewShortLinkService(shortLinkRepo)
	shortLinkCtrl := rest.NewShortLinkController(shortLinkService)

	authRepo := repository.NewUserRepository()
	authService := service.NewAuthService(authRepo, jwtAuth)
	authCtrl := rest.NewAuthController(authService)

	router.HandleFunc("/ping/", handlePing).Methods("GET")
	router.HandleFunc("/s/{shortName}", shortLinkCtrl.RedirectHandler)

	// v1.HandleFunc("/shorten", shortLinkCtrl.ShortenHandler)
	// v1.HandleFunc("/edit/{shortName}", shortLinkCtrl.EditShortLinkHandler)
	// v1.HandleFunc("/delete/{shortName}", shortLinkCtrl.DeleteShortLinkHandler)

	v1.HandleFunc("/auth/login/", authCtrl.LoginUser)
	v1.HandleFunc("/auth/register/", authCtrl.RegisterUser)

	v1.HandleFunc("/shorten", apply(shortLinkCtrl.ShortenHandler, requireAuth))
	v1.HandleFunc("/edit/{shortName}", apply(shortLinkCtrl.EditShortLinkHandler, requireAuth))
	v1.HandleFunc("/delete/{shortName}", apply(shortLinkCtrl.DeleteShortLinkHandler, requireAuth))

	err := http.ListenAndServe(":4321", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
