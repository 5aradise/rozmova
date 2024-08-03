package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createJWTtoken(id int, expInSec int64) (string, error) {
	var dayInSec int64 = 24 * 60 * 60
	if expInSec <= 0 || expInSec > dayInSec {
		expInSec = dayInSec
	}
	now := time.Now().UTC()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "rozmova",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expInSec * 1000000000))),
			Subject:   strconv.Itoa(id),
		})
	return t.SignedString(cfg.jwtKey)
}

func (cfg *apiConfig) getJWTtoken(r *http.Request) (*jwt.Token, error) {
	authArr := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authArr) != 2 || authArr[0] != "Bearer" {
		return nil, errors.New("wrong authorization field")
	}

	jwtToken, err := jwt.ParseWithClaims(authArr[1], &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return cfg.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}
