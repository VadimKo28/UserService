package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewTokenService(jwtSecret []byte, tokenTTL time.Duration) *TokenService {
	return &TokenService{jwtSecret: jwtSecret, tokenTTL: tokenTTL}
}

func (s *TokenService) SignToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
	})

	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) ParseToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return "", err
	}

	if !token.Valid || claims.Subject == "" {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}
