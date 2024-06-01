package main

import (
	"encoding/json"
	"net/http"

	"nhuxoll/bdd_chirpy/internal/auth"
)

type ReturnUser struct {
	ID           int    `json:"id"`
	EMail        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
}
type User struct {
	ID       int    `json:"id"`
	EMail    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
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
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
	}
	user, err := cfg.DB.CreateUser(params.EMail, hashedPass)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User already exists")
		return
	}
	respondWithJSON(w, http.StatusCreated, ReturnUser{
		ID:    user.ID,
		EMail: user.EMail,
	})
}

// func validateUser(body string) (string, error) {
// 	const maxUserLength = 140
// 	if len(body) > maxUserLength {
// 		return "", errors.New("User is too long")
// 	}
//
// 	badWords := map[string]struct{}{
// 		"kerfuffle": {},
// 		"sharbert":  {},
// 		"fornax":    {},
// 	}
// 	cleaned := getCleanedBody(body, badWords)
// 	return cleaned, nil
// }

// func getCleanedBody(body string, badWords map[string]struct{}) string {
// 	words := strings.Split(body, " ")
// 	for i, word := range words {
// 		loweredWord := strings.ToLower(word)
// 		if _, ok := badWords[loweredWord]; ok {
// 			words[i] = "****"
// 		}
// 	}
// 	cleaned := strings.Join(words, " ")
// 	return cleaned
// }
