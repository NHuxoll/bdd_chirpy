package main

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
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

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't interpred chirpID")
	}
	userIDStr := claims.Subject
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID in token")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID, userID)
	if err != nil {
		respondWithError(w, 403, "User not author of this chirp")
		return
	}

	respondWithJSON(w, 204, Chirp{})
}
