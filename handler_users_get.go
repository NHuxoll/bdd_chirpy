package main

import (
	"net/http"
	// "sort"
	"strconv"
)

// func (cfg *apiConfig) handlerUsersRetrieveAll(w http.ResponseWriter, r *http.Request) {
// 	dbUsers, err := cfg.DB.GetUsers()
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user")
// 		return
// 	}
//
// 	users := []User{}
// 	for _, dbUser := range dbUsers {
// 		users = append(users, User{
// 			ID:    dbUser.ID,
// 			EMail: dbUser.EMail,
// 		})
// 	}
//
// 	sort.Slice(users, func(i, j int) bool {
// 		return users[i].ID < users[j].ID
// 	})
//
// 	respondWithJSON(w, http.StatusOK, users)
// }

func (cfg *apiConfig) handlerUsersRetrieveById(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't interpred userID")
	}

	user, err := cfg.DB.GetUserById(userID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve user")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
