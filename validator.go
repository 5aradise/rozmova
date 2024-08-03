package main

import (
	"errors"
	"regexp"
)

var badWordsRegex = regexpMustCompileSlice([]string{"[ХхXx]{1,}[0ОоOo]{1,}[ХхXx]{1,}[06ОоOoАаAa]{0,}[!1ЛлLl]{1,}[!1ИиІіЫыSsIi]{0,}",
	"порохобот",
	"зєлєбобік"})

func regexpMustCompileSlice(exprs []string) []*regexp.Regexp {
	regexps := make([]*regexp.Regexp, 0, len(exprs))
	for _, expr := range exprs {
		regexps = append(regexps, regexp.MustCompile(expr))
	}
	return regexps
}

func validateMessage(msg string) (string, error) {
	const validLen = 140
	if len(msg) > validLen {
		return "", errors.New("message is too long")
	}
	return cleanMessage(msg), nil
}

func cleanMessage(msg string) string {
	const replaceStr = "****"
	return replaceBadWords(msg, replaceStr)
}

func replaceBadWords(msg, replace string) string {
	for _, wordRegex := range badWordsRegex {
		msg = wordRegex.ReplaceAllString(msg, replace)
	}
	return msg
}
