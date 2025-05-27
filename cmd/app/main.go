package main

import (
	"encoding/json"
	"javaneseivankov/url-shortener/internal/middlewares"
	"javaneseivankov/url-shortener/internal/repostiories"
	"javaneseivankov/url-shortener/pkg"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type RequestShortenURL struct {
	Url string `json:"url"`
	ShortName string `json:"short_name"`
}

type ResponseShortenURL struct {
	Url string `json:"url"`
}

var repo = repostiories.NewShortLinkRepository()


func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req RequestShortenURL;

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error: Bad Request", http.StatusBadRequest)
	}

	repoRes , err := repo.CreateRedirectLink(req.ShortName, req.Url);
	redirectUrl :=  "/s/" + repoRes.ShortName

	 if err != nil {
		http.Error(w, "Error: Internal Server Error", http.StatusInternalServerError) 
	} 

	res := RequestShortenURL{
		Url: redirectUrl,
	}

	pkg.RenderJSON(w, res)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortName := vars["shortName"]

	repoRes, err := repo.GetRedirectLink(shortName);
	redirectUrl := repoRes.OriginalUrl

	if err != nil {
		http.Error(w, "Error: Internal Server Error", http.StatusInternalServerError) 
	}

	http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
}

func apiPathV1(path string) string {
	return "/api/v1" + path 
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc(apiPathV1("/shorten"), shortenHandler).Methods("POST")
	router.HandleFunc("/s/{shortName}", redirectHandler).Methods("GET")
	router.Use(middlewares.LoggingMiddleware)

	err := http.ListenAndServe(":4321", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
