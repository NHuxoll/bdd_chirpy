package database

import (
	"errors"
)

type Chirp struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author_ID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:        id,
		Body:      body,
		Author_ID: author_id,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if chirp, ok := dbStructure.Chirps[id]; !ok {

		return Chirp{}, errors.New("Chirp not found")
	} else {

		return chirp, nil
	}
}

func (db *DB) DeleteChirp(chirpId, userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	chirp, ok := dbStructure.Chirps[chirpId]
	if !ok {
		return errors.New("Chirp not found")
	}
	if chirp.Author_ID != userId {
		return errors.New("User not creator!")
	}
	delete(dbStructure.Chirps, chirpId)

	return nil
}
