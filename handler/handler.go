package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func HandlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleandedBody := cleanBody(params.Body)
	RespondWithJSON(w, http.StatusOK, returnVals{
		Cleaned: cleandedBody,
	})
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func cleanBody(body string) string {
	cleanBody := []string{}
	for _, word := range strings.Split(body, " ") {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			cleanBody = append(cleanBody, "****")
		default:
			cleanBody = append(cleanBody, word)
		}
	}

	return strings.Join(cleanBody, " ")

}
