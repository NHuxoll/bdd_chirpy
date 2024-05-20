package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	EMail string `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		EMail string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.EMail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
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

