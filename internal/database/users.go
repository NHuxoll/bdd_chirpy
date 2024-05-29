package database

import (
	"errors"
)

type User struct {
	Password string `json:"password"`
	ID       int    `json:"id"`
	EMail    string `json:"email"`
}

var ErrAlreadyExists = errors.New("already exists")
var ErrNotExist = errors.New("does not exist")

func (db *DB) CreateUser(email string, hashedPass string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		ID:       id,
		Password: hashedPass,
		EMail:    email,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// func (db *DB) GetUsers() ([]User, error) {
// 	dbStructure, err := db.loadDB()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	users := make([]User, 0, len(dbStructure.Users))
// 	for _, user := range dbStructure.Users {
// 		users = append(users, user)
// 	}
//
// 	return users, nil
// }

func (db *DB) UpdateUser(email, password string, id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err

	}
	user, ok := dbStructure.Users[id]
	if !ok {

		return User{}, errors.New("User not found")
	}
	if password != "" {
		user.Password = password
	}
	if email != "" {
		user.EMail = email
	}
	dbStructure.Users[id] = user

	if err = db.writeDB(dbStructure); err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	if user, ok := dbStructure.Users[id]; !ok {

		return User{}, errors.New("User not found")
	} else {
		return user, nil
	}
}

func (db *DB) GetUserByEmail(mail string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if mail == user.EMail {
			return user, nil
		}
	}
	return User{}, ErrNotExist
}
