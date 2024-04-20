package database

import (
	"strconv"

	"github.com/5aradise/jsondb"
)

type Message struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

var msgPath = "messages"

func (db *DB) AddMsg(data string) (int, error) {
	id, err := db.GetLen(msgPath)
	if err != nil {
		return 0, err
	}

	msg := Message{
		Id:   id,
		Data: data,
	}
	err = db.Insert(msgPath+db.Divider()+strconv.Itoa(id), msg)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) ReadMsgById(id string) (Message, error) {
	msg := Message{}
	err := db.GetStruct(msgPath+db.Divider()+id, &msg)
	if err != nil {
		return Message{}, err
	}
	return msg, nil
}

func (db *DB) ReadMsgs() ([]Message, error) {
	maps, err := db.GetAllMaps(msgPath)
	if err != nil {
		return nil, err
	}

	msgs := make([]Message, 0, len(maps))
	for _, mapInst := range maps {
		msg := Message{}
		err = jsondb.MapToStruct(&msg, mapInst)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}
