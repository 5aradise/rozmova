package database

import (
	"errors"
	"strconv"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) AddUser(email string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return err
	}

	id := len(dbStruct.Users) + 1
	dbStruct.Users[strconv.Itoa(id)] = User{
		Id:    id,
		Email: email,
	}
	err = db.writeDB(dbStruct)
	return err
}

func (db *DB) ReadUser(id string) (User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStruct.Users[id]
	if !ok {
		return User{}, errors.New("user with this id doesnt exist")
	}

	return user, nil
}

func (db *DB) ReadUsers() ([]User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStruct.Users))
	for _, user := range dbStruct.Users {
		users = append(users, user)
	}
	return users, nil
}
