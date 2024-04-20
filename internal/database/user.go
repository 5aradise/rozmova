package database

import (
	"errors"
	"strconv"
)

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"hashedPassword"`
}

func (db *DB) AddUser(email string, hashedPassword []byte) (int, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return 0, err
	}

	id := len(dbStruct.Users) + 1
	dbStruct.Users[strconv.Itoa(id)] = User{
		Id:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	err = db.writeDB(dbStruct)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) UpdateUser(id string, email string, hashedPassword []byte) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return User{}, err
	}

	updatedUser, ok := dbStruct.Users[id]
	if !ok {
		return User{}, errors.New("user with this id doesnt exist")
	}

	if email != "" {
		updatedUser.Email = email
	}
	if len(hashedPassword) != 0 {
		updatedUser.HashedPassword = hashedPassword
	}

	dbStruct.Users[id] = updatedUser

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (db *DB) ReadUserById(id string) (User, error) {
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

func (db *DB) ReadUserByEmail(email string) (User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user with this email doesnt exist")
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
