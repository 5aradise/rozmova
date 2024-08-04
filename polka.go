package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) polkaWebhooks(w http.ResponseWriter, r *http.Request) {
	authArr := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authArr) != 2 || authArr[0] != "ApiKey" {
		respondWithError(w, http.StatusUnauthorized, "wrong authorization field")
		return
	}
	key := authArr[1]
	if key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "wrong key")
		return
	}

	type reqWebhook struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	webhook := reqWebhook{}
	decodeReq(r, &webhook)

	if webhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
		return
	}

	_, err := cfg.db.UpdateUserSubscription(webhook.Data.UserId,true)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}
