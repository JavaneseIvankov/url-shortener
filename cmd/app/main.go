package main

import (
	"javaneseivankov/url-shortener/internal/controller/rest"
	"javaneseivankov/url-shortener/internal/middleware"
	"javaneseivankov/url-shortener/internal/repository"
	"javaneseivankov/url-shortener/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


func main() {
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)

	shortLinkRepo := repository.NewShortLinkRepository()
	shortLinkService := service.NewShortLinkService(shortLinkRepo)
	shortLinkCtrl := rest.NewShortLinkController(shortLinkService)

	router.HandleFunc(rest.ApiPathV1("/shorten"), shortLinkCtrl.ShortenHandler)
	router.HandleFunc(rest.ApiPathV1("/edit/{shortName}"), shortLinkCtrl.EditShortLinkHandler)
	router.HandleFunc("/s/{shortName}", shortLinkCtrl.RedirectHandler)

	err := http.ListenAndServe(":4321", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
