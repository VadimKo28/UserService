package service

import (
	"crypto/rand"
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

func (s *TokenService) GenerateTokens(userID int) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
	})

	accessToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
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

func newRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
