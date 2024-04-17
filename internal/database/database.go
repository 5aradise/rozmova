package database

import (
	db "github.com/5aradise/jsondb"
)

type DB struct {
	db.Jsondb
}

func NewDB(path string, divider ...string) (*DB, error) {
	db, err := db.New(path, divider...)
	if err != nil {
		return nil, err
	}

	databse := &DB{
		*db,
	}

	databse.InsertDir(msgPath)
	databse.InsertDir(userPath)

	return databse, nil
}
