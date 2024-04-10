package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStruct struct {
	Messages map[string]Message `json:"messages"`
	Users    map[string]User    `json:"users"`
}

func NewDB(path string) (*DB, error) {
	if !isFileExist(path) {
		data, err := json.Marshal(DBStruct{
			Messages: make(map[string]Message),
			Users:    make(map[string]User)})
		if err != nil {
			return &DB{}, err
		}

		err = os.WriteFile(path, data, 0600)
		if err != nil {
			return &DB{}, err
		}
	}

	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}, nil
}

func (db *DB) readDB() (DBStruct, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStruct{}, err
	}

	readedDB := DBStruct{}
	err = json.Unmarshal(data, &readedDB)
	if err != nil {
		return DBStruct{}, err
	}

	return readedDB, nil
}

func (db *DB) writeDB(dbToWrite DBStruct) error {
	data, err := json.Marshal(dbToWrite)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0600)
	return err
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
