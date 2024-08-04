package database

import (
	"errors"
	"strconv"
	"time"

	"github.com/5aradise/jsondb"
)

type User struct {
	Id             int          `json:"id"`
	Email          string       `json:"email"`
	HashedPassword []byte       `json:"hashedPassword"`
	IsSub          bool         `json:"isSub"`
	RefreshToken   RefreshToken `json:"refreshToken"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expAt"`
}

var userPath = "users"

func (db *DB) AddUser(email string, hashedPassword []byte) (User, error) {
	_, err := db.ReadUserByEmail(email)
	if err == nil {
		return User{}, errors.New("user with this email already registered")
	}

	id, err := db.GetLen(userPath)
	if err != nil {
		return User{}, err
	}
	id++

	user := User{
		Id:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	err = db.Insert(userPath+db.Divider()+strconv.Itoa(id), user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUserEmail(id int, email string) (User, error) {
	currUserPath := userPath + db.Divider() + strconv.Itoa(id)

	updatedUser, err := db.ReadUserById(id)
	if err != nil {
		return User{}, errors.New("user with this id doesnt exist")
	}

	err = db.Insert(currUserPath+db.Divider()+"email", email)
	if err != nil {
		return User{}, err
	}
	updatedUser.Email = email

	return updatedUser, nil
}

func (db *DB) UpdateUserPassword(id int, hashedPassword []byte) (User, error) {
	currUserPath := userPath + db.Divider() + strconv.Itoa(id)

	updatedUser, err := db.ReadUserById(id)
	if err != nil {
		return User{}, errors.New("user with this id doesnt exist")
	}

	err = db.Insert(currUserPath+db.Divider()+"hashedPassword", hashedPassword)
	if err != nil {
		return User{}, err
	}
	updatedUser.HashedPassword = hashedPassword

	return updatedUser, nil
}

func (db *DB) UpdateUserToken(id int, token string) (User, error) {
	currUserPath := userPath + db.Divider() + strconv.Itoa(id)

	updatedUser, err := db.ReadUserById(id)
	if err != nil {
		return User{}, errors.New("user with this id doesnt exist")
	}

	const expTime = time.Hour * 24 * 60
	err = db.Insert(currUserPath+db.Divider()+"refreshToken", RefreshToken{
		Token:     token,
		ExpiresAt: time.Now().Add(expTime),
	})
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (db *DB) UpdateUserSubscription(id int, isSub bool) (User, error) {
	currUserPath := userPath + db.Divider() + strconv.Itoa(id)

	updatedUser, err := db.ReadUserById(id)
	if err != nil {
		return User{}, errors.New("user with this id doesnt exist")
	}

	err = db.Insert(currUserPath+db.Divider()+"isSub", isSub)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (db *DB) ReadUserById(id int) (User, error) {
	user := User{}
	err := db.GetStruct(userPath+db.Divider()+strconv.Itoa(id), &user)
	if err != nil {
		return User{}, errors.New("user with this id doesnt exist")
	}
	return user, nil
}

func (db *DB) ReadUserByEmail(email string) (User, error) {
	for i := 1; ; i++ {
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

func (db *DB) ReadUserByToken(token string) (User, error) {
	user := User{}
	for i := 1; ; i++ {
		err := db.GetStruct(userPath+db.Divider()+strconv.Itoa(i), &user)
		if err != nil {
			return User{}, errors.New("user with this token doesnt exist")
		}

		if user.RefreshToken.Token == token {
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
