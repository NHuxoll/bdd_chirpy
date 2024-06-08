package main

import (
	"encoding/json"
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
)

type Chirp struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author_ID int    `json:"author_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	header := r.Header.Get("Authorization")

	// Make sure the provided header is not empty or malformed
	if header == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing or malformed token")
		return
	}

	// Extract the token part from "Bearer <token>"
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		respondWithError(w, http.StatusUnauthorized, "Malformed token")
		return
	}

	tokenString := parts[1]
	claims := &jwt.RegisteredClaims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userIDStr := claims.Subject
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID in token")
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		Author_ID: chirp.Author_ID,
		Body:      chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
