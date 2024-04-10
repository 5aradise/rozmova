package main

import (
	"encoding/json"
	"net/http"
	"sort"
)

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type respUser struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	user := respUser{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}

	if len(user.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "empty email")
		return
	}

	err = cfg.db.AddUser(user.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"email": user.Email})
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
	user, err := cfg.db.ReadUser(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
