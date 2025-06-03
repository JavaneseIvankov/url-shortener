package main

import (
	"javaneseivankov/url-shortener/internal/controller/rest"
	"javaneseivankov/url-shortener/internal/middleware"
	"javaneseivankov/url-shortener/internal/repository"
	"javaneseivankov/url-shortener/internal/service"
	"javaneseivankov/url-shortener/pkg"
	"javaneseivankov/url-shortener/pkg/db"
	"javaneseivankov/url-shortener/pkg/jwt"
	"javaneseivankov/url-shortener/pkg/logger"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	pkg.SendJSON(w, http.StatusOK, map[string]interface{}{"payload": "pong"})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	environ := os.Getenv("ENV")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtTTL := os.Getenv("JWT_TTL")

	if environ == "" {
		log.Fatalln("ENV must be set")
	}

	if jwtSecret == "" || jwtTTL == ""  {
		log.Fatalln("JWT_SECRET and JWT_TTL must be set")
	}

	// TODO: Refactor to avoid error, temporal coupling
	logger.Init(environ)
	jwtAuth := jwt.NewJWT(jwtSecret, jwtTTL)

	if err := db.Init(); err != nil {
		panic("Failed to initialize DB: " + err.Error())
	}
	
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

	err = http.ListenAndServe(":4321", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
