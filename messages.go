package main

import (
	"net/http"
	"sort"
)

func (cfg *apiConfig) postMessage(w http.ResponseWriter, r *http.Request) {
	type respMessage struct {
		Body string `json:"body"`
	}

	msg := respMessage{}
	err := decodeResp(r, &msg)
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

	id, err := cfg.db.AddMsg(cleanMsg)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]any{"id": id, "body": cleanMsg})
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
	respondWithJSON(w, http.StatusOK, msgs)
}

func (cfg *apiConfig) getMessageById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("messageId")
	msg, err := cfg.db.ReadMsgById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, msg)
}
