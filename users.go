package main

import (
	"net/http"
	"sort"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) registerUser(w http.ResponseWriter, r *http.Request) {
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

	if len(req.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "empty email")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.db.AddUser(req.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) changeUser(w http.ResponseWriter, r *http.Request) {
	userId, err := cfg.getIdFromJWT(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type reqUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := reqUser{}
	err = decodeReq(r, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "empty password")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Email != "" {
		_, err = cfg.db.UpdateUserEmail(userId, req.Email)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	user, err := cfg.db.UpdateUserPassword(userId, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
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
