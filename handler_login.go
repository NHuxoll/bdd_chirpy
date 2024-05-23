package main

import (
	"encoding/json"
	"net/http"
	"nhuxoll/bdd_chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		EMail    string `json:"email"`
	}

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
	respondWithJSON(w, http.StatusOK, ReturnUser{ID: user.ID, EMail: user.EMail})
}
