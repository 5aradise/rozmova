package database

import (
	"errors"
	"strconv"

	"github.com/5aradise/jsondb"
)

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"hashedPassword"`
}

var userPath = "users"

func (db *DB) AddUser(email string, hashedPassword []byte) (int, error) {
	id, err := db.GetLen(userPath)
	if err != nil {
		return 0, err
	}

	user := User{
		Id:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	err = db.Insert(userPath+db.Divider()+strconv.Itoa(id), user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) ReadUserById(id string) (User, error) {
	user := User{}
	err := db.GetStruct(userPath+db.Divider()+id, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) ReadUserByEmail(email string) (User, error) {
	for i := 0; ; i++ {
		mappedUser, err := db.GetMap(userPath + db.Divider() + strconv.Itoa(i))
		if err != nil {
			return User{}, errors.New("user with this email doesnt exist")
		}

		if mappedUser["email"] == email {
			user := User{}
			err = jsondb.MapToStruct(&user, mappedUser)
			if err != nil {
				return User{}, err
			}
			return user, nil
		}
	}
}

func (db *DB) ReadUsers() ([]User, error) {
	maps, err := db.GetAllMaps(userPath)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(maps))
	for _, mapInst := range maps {
		user := User{}
		err = jsondb.MapToStruct(&user, mapInst)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
