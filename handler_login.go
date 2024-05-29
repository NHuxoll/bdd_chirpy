package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nhuxoll/bdd_chirpy/internal/auth"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password   string `json:"password"`
		EMail      string `json:"email"`
		ExpireTime int    `json:"expires_in_seconds"`
	}
	var expTime time.Duration
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	user, err := cfg.DB.GetUserByEmail(params.EMail)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve user")
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return

	}

	if params.ExpireTime == 0 {
		expTime = time.Duration(time.Hour * 24)
	} else {
		expTime = time.Duration(int(time.Second) * params.ExpireTime)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expTime)),
			Subject:   strconv.Itoa(user.ID)})
	jwt, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	respondWithJSON(w, http.StatusOK, ReturnUser{ID: user.ID, EMail: user.EMail, Token: jwt})
}
