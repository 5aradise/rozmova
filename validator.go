package main

import (
	"errors"
	"regexp"
)

func validateMessage(msg string) (string, error) {
	const validLen = 140
	if len(msg) > validLen {
		return "", errors.New("message is too long")
	}
	return cleanMessage(msg), nil
}

func cleanMessage(msg string) string {
	const replaceStr = "ІДІ НАХУЙ"
	return replaceBadWords(msg, replaceStr)
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
