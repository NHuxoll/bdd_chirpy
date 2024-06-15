package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"nhuxoll/bdd_chirpy/internal/database"
	"strings"
	// "nhuxoll/bdd_chirpy/internal/auth"
	// "time"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {

	apiHeader := r.Header.Get("Authorization")
	if apiHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "API key not found")
		return
	}
	splitAPI := strings.Split(apiHeader, " ")
	if len(splitAPI) < 2 || splitAPI[0] != "ApiKey" {
		respondWithError(w, http.StatusUnauthorized, "API key malformed")
	}

	if splitAPI[1] != cfg.polkaKey {
		w.WriteHeader(401)
		return
	}
	type data struct {
		UserID int `json:"user_id"`
	}
	// Token is valid, proceed with your logic to extract the user ID, etc.

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Check if the user exists in the database
	_, err = cfg.DB.UpgradeUserStatus(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Couldn't find the user")
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update the user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
