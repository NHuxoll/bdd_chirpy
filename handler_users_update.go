package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"nhuxoll/bdd_chirpy/internal/auth"

	jwt "github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		EMail    string `json:"email"`
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

	// Token is valid, proceed with your logic to extract the user ID, etc.

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
	}
	// Extract user ID from token claims
	userIDStr := claims.Subject
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID in token")
		return
	}

	// Check if the user exists in the database
	user, err := cfg.DB.GetUserById(userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User does not exist")
		return
	}

	user, err = cfg.DB.UpdateUser(params.EMail, hashedPass, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User already exists")
		return
	}
	respondWithJSON(w, http.StatusCreated, ReturnUser{
		ID:    user.ID,
		EMail: user.EMail,
	})
}
