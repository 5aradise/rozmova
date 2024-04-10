package main

import (
	"net/http"
	"sort"

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

	_, err = cfg.db.ReadUserByEmail(resp.Email)
	if err == nil {
		respondWithError(w, http.StatusUnauthorized, "user with this email already registered")
		return
	}
	if err.Error() != "user with this email doesnt exist" {
		respondWithError(w, http.StatusBadRequest, err.Error())
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
	id := r.PathValue("userId")
	user, err := cfg.db.ReadUserById(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
