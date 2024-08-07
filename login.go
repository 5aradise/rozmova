package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type reqUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req := reqUser{}
	err := decodeReq(r, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}

	requiredUser, err := cfg.db.ReadUserByEmail(req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong password")
		return
	}

	err = bcrypt.CompareHashAndPassword(requiredUser.HashedPassword, []byte(req.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong password")
		return
	}

	accessToken, err := cfg.createJWTtoken(requiredUser.Id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	refreshToken, err := createRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = cfg.db.UpdateUserToken(requiredUser.Id, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"id":            requiredUser.Id,
		"email":         requiredUser.Email,
		"token":         accessToken,
		"refresh_token": refreshToken,
		"is_chirpy_red": requiredUser.IsSub,
	})
}
