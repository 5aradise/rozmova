package database

import (
	"errors"
	"strconv"
)

type Message struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) AddMsg(body string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return err
	}

	id := len(dbStruct.Messages) + 1
	dbStruct.Messages[strconv.Itoa(id)] = Message{
		Id:   id,
		Body: body,
	}
	err = db.writeDB(dbStruct)
	return err
}

func (db *DB) ReadMsg(id string) (Message, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return Message{}, err
	}

	message, ok := dbStruct.Messages[id]
	if !ok {
		return Message{}, errors.New("message with this id doesnt exist")
	}

	return message, nil
}

func (db *DB) ReadMsgs() ([]Message, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStruct, err := db.readDB()
	if err != nil {
		return nil, err
	}

	messages := make([]Message, 0, len(dbStruct.Messages))
	for _, msg := range dbStruct.Messages {
		messages = append(messages, msg)
	}
	return messages, nil
}
