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

// START DTO
type RequestShortenURL struct {
	Url string `json:"url"`
	ShortName string `json:"short_name"`
}

type ResponseShortenURL struct {
	Url string `json:"url"`
}

type RequestEditShortLink struct {
	NewUrl string `json:"new_url"`
}

type ResponseEditShortLink struct {
	Url string `json:"new_url"`
	ShortName string `json:"short_name"`
}
// END DTO

// START INITIALIZER
var repo = repostiories.NewShortLinkRepository()
// END INITIALIZER


// START HANDLER
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req RequestShortenURL;

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.BadRequestErrorDefault(w)
		return
	}

	// START SERVICE LAYER
	repoRes , err := repo.CreateRedirectLink(req.ShortName, req.Url);
	redirectUrl :=  "/s/" + repoRes.ShortName

	 if err != nil {
		pkg.InternalServerErrorDefault(w)
		return
	} 

	res := ResponseShortenURL{
		Url: redirectUrl,
	}
	// END SERVICE LAYER

	pkg.RenderJSON(w, res)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortName := vars["shortName"]

	// START SERVICE LAYER
	repoRes, err := repo.GetRedirectLink(shortName);
	redirectUrl := repoRes.OriginalUrl

	if err != nil {
		http.Error(w, "Error: Internal Server Error", http.StatusInternalServerError) 
		return
	}
	// END SERVICE LAYER

	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func editShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortName := vars["shortName"]

	var req RequestEditShortLink

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.BadRequestErrorDefault(w)
		return
	}

	// START SERVICE LAYER
	repoRes, err := repo.EditShortLink(shortName, req.NewUrl); if err != nil {
		pkg.InternalServerError(w, err.Error())
		return
	}
	redirectUrl :=  "/s/" + repoRes.ShortName

	res := &ResponseEditShortLink{
		Url: redirectUrl,
		ShortName: repoRes.ShortName,
	}
	// END SERVICE LAYER
	
	pkg.RenderJSON(w, res)
}

func apiPathV1(path string) string {
	return "/api/v1" + path 
}

// END HANDLER

func main() {
	router := mux.NewRouter()
	router.HandleFunc(apiPathV1("/shorten"), shortenHandler).Methods("POST")
	router.HandleFunc(apiPathV1("/edit/{shortName}"), editShortLinkHandler).Methods("POST")
	router.HandleFunc("/s/{shortName}", redirectHandler).Methods("GET")
	router.Use(middlewares.LoggingMiddleware)

	err := http.ListenAndServe(":4321", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
