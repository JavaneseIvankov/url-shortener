package pkg

import "net/http"

func InternalServerError(w http.ResponseWriter, errorMsg string) {
	http.Error(w, errorMsg, http.StatusInternalServerError)
}

func InternalServerErrorDefault(w http.ResponseWriter) {
	http.Error(w, "Error: Internal Server Error", http.StatusInternalServerError)
}

func BadRequestError(w http.ResponseWriter, errorMsg string) {
		http.Error(w, errorMsg, http.StatusBadRequest)
}

func BadRequestErrorDefault(w http.ResponseWriter) {
		BadRequestError(w,  "Error: Bad Request") 
}
