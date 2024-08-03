package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	authArr := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authArr) != 2 || authArr[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "1")
		return
	}

	refreshToken := authArr[1]
	user, err := cfg.db.ReadUserByToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "2")
		return
	}

	if user.RefreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "3")
		return
	}

	newToken, err := cfg.createJWTtoken(user.Id)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "4")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newToken})
}

func createRefreshToken() (string, error) {
	byteToken := make([]byte, 32)
	_, err := rand.Read(byteToken)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(byteToken), nil
}
