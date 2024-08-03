package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	authArr := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authArr) != 2 || authArr[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "wrong authorization field")
		return
	}

	refreshToken := authArr[1]
	user, err := cfg.db.ReadUserByToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong authorization field")
		return
	}

	randToken, err := createRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = cfg.db.UpdateUser(user.Id, "", nil, randToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}
