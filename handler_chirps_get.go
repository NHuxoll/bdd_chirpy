package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsRetrieveAll(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	sortOrder := "asc"

	if r.URL.Query().Get("sort") == "desc" {
		sortOrder = "desc"
	}

	if r.URL.Query().Has("author_id") && r.URL.Query().Get("author_id") != "" {
		id := r.URL.Query().Get("author_id")
		authorID, err := strconv.Atoi(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse author id")
			return
		}
		chirps, err := cfg.DB.GetChirpByAuthorId(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve the chirps")
		}
		if len(chirps) > 0 {

			sort.Slice(chirps, func(i, j int) bool {
				if sortOrder == "asc" {
					return chirps[i].ID < chirps[j].ID
				} else {

					return chirps[i].ID > chirps[j].ID
				}
			})
			respondWithJSON(w, http.StatusOK, chirps)
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			Body:      dbChirp.Body,
			Author_ID: dbChirp.Author_ID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "asc" {
			return chirps[i].ID < chirps[j].ID
		} else {

			return chirps[i].ID > chirps[j].ID
		}
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsRetrieveById(w http.ResponseWriter, r *http.Request) {

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't interpred chirpID")
	}

	chirp, err := cfg.DB.GetChirpById(chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
