package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type respUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expires  int    `json:"expires_in_seconds"`
	}

	resp := respUser{}
	err := decodeResp(r, &resp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}

	requiredUser, err := cfg.db.ReadUserByEmail(resp.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong password")
		return
	}

	err = bcrypt.CompareHashAndPassword(requiredUser.HashedPassword, []byte(resp.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong password")
		return
	}

	jwtToken, err := cfg.createJWTtoken(requiredUser.Id, int64(resp.Expires))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"token": jwtToken,
	})
}
