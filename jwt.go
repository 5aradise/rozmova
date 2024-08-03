package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createJWTtoken(userId int) (string, error) {
	const expTime = time.Hour
	now := time.Now().UTC()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "rozmova",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expTime)),
			Subject:   strconv.Itoa(userId),
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

func (cfg *apiConfig) getIdFromJWT(r *http.Request) (int, error) {
	token, err := cfg.getJWTtoken(r)
	if err != nil {
		return 0, err
	}

	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
