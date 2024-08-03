package main

import (
	"net/http"
	"sort"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) registerUser(w http.ResponseWriter, r *http.Request) {
	type respUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	resp := respUser{}
	err := decodeResp(r, &resp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}

	if len(resp.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "empty email")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resp.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := cfg.db.AddUser(resp.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]any{"email": resp.Email, "id": id})
}

func (cfg *apiConfig) changeUser(w http.ResponseWriter, r *http.Request) {
	userId, err := cfg.getIdFromJWT(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type respUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	resp := respUser{}
	err = decodeResp(r, &resp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if resp.Password == "" {
		respondWithError(w, http.StatusBadRequest, "empty password")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resp.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.db.UpdateUser(userId, resp.Email, hashedPassword, "")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"email": user.Email, "id": user.Id})
}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.ReadUsers()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})
	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.db.ReadUserById(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
