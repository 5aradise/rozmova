package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func messageValidator(w http.ResponseWriter, r *http.Request) {
	const validLen = 140
	type message struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	msg := message{}
	err := decoder.Decode(&msg)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	if len(msg.Body) > validLen {
		respondWithError(w, http.StatusBadRequest, "Message is too long")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": cleanMessage(msg.Body)})
}

func cleanMessage(msg string) string {
	return replaceBadWords(msg, "ІДІ НАХУЙ")
}

func replaceBadWords(msg, replace string) string {
	badWordsExp := []string{"[ХхXx]{1,}[0ОоOo]{1,}[ХхXx]{1,}[06ОоOoАаAa]{0,}[!1ЛлLl]{1,}[!1ИиІіЫыSsIi]{0,}", "порохобот", "зєлєбобік"}
	for _, wordExp := range badWordsExp {
		wordRegex, err := regexp.Compile(wordExp)
		if err != nil {
			continue
		}
		msg = wordRegex.ReplaceAllString(msg, replace)
	}
	return msg
}
