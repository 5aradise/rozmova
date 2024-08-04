package database

import (
	"strconv"

	"github.com/5aradise/jsondb"
)

type Message struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

var msgPath = "messages"

func (db *DB) AddMsg(authorId int, data string) (int, error) {
	id, err := db.GetLen(msgPath)
	if err != nil {
		return 0, err
	}
	id++

	msg := Message{
		Id:       id,
		Body:     data,
		AuthorId: authorId,
	}
	err = db.Insert(msgPath+db.Divider()+strconv.Itoa(id), msg)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) ReadMsgById(id int) (Message, error) {
	msg := Message{}
	err := db.GetStruct(msgPath+db.Divider()+strconv.Itoa(id), &msg)
	if err != nil {
		return Message{}, err
	}
	return msg, nil
}

func (db *DB) ReadMsgs() ([]*Message, error) {
	maps, err := db.GetAllMaps(msgPath)
	if err != nil {
		return nil, err
	}

	msgs := make([]*Message, 0, len(maps))
	for _, mapInst := range maps {
		msg := &Message{}
		err = jsondb.MapToStruct(msg, mapInst)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (db *DB) DeleteMsg(id int) error {
	return db.Delete(msgPath + db.Divider() + strconv.Itoa(id))
}
