package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/5aradise/rozmova/internal/database"
)

func (cfg *apiConfig) postMessage(w http.ResponseWriter, r *http.Request) {
	type reqMessage struct {
		Body string `json:"body"`
	}

	userId, err := cfg.getIdFromJWT(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong jwt token")
		return
	}

	msg := reqMessage{}
	err = decodeReq(r, &msg)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(msg.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "empty message")
		return
	}

	cleanMsg, err := validateMessage(msg.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	msgId, err := cfg.db.AddMsg(userId, cleanMsg)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]any{"id": msgId, "body": cleanMsg, "author_id": userId})
}

func (cfg *apiConfig) getMessages(w http.ResponseWriter, r *http.Request) {
	msgs, err := cfg.db.ReadMsgs()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Id < msgs[j].Id
	})

	authorIdStr := r.URL.Query().Get("author_id")
	if authorIdStr != "" {
		authorId, err := strconv.Atoi(authorIdStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		authorsMsgs := make([]*database.Message, 0)
		for _, msg := range msgs {
			if msg.AuthorId == authorId {
				authorsMsgs = append(authorsMsgs, msg)
			}
		}
		respondWithJSON(w, http.StatusOK, authorsMsgs)
		return
	}

	respondWithJSON(w, http.StatusOK, msgs)
}

func (cfg *apiConfig) getMessageById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("messageId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg, err := cfg.db.ReadMsgById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, msg)
}

func (cfg *apiConfig) deleteMessage(w http.ResponseWriter, r *http.Request) {
	userId, err := cfg.getIdFromJWT(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong jwt token")
		return
	}

	msgId, err := strconv.Atoi(r.PathValue("messageId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg, err := cfg.db.ReadMsgById(msgId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if msg.AuthorId != userId {
		respondWithError(w, http.StatusForbidden, "wrong jwt token")
		return
	}

	err = cfg.db.DeleteMsg(msgId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}
