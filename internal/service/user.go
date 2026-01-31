package service

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"user_advt/internal/domain/users"

	"github.com/golang-jwt/jwt/v5"
)

type UserStorage interface {
	Save(ctx context.Context, user *users.UserCreateDTO) (int, error)
	Get(ctx context.Context, id string) (users.User, error)
	GetByCredentials(ctx context.Context, user *users.UserSignInDTO) (users.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserService struct {
	storage   UserStorage
	hasher    PasswordHasher
	tokenTTL  time.Duration
	jwtSecret []byte
}

func NewUserService(storage UserStorage, hasher PasswordHasher, ttl time.Duration, jwtSecret []byte) *UserService {
	return &UserService{storage: storage, hasher: hasher, tokenTTL: ttl, jwtSecret: jwtSecret}
}

func (s *UserService) SignUpUser(ctx context.Context, name, email, password string) (int, error) {
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return 0, err
	}

	user := users.UserCreateDTO{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	return s.storage.Save(ctx, &user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (users.User, error) {

	user, err := s.storage.Get(ctx, id)

	if err != nil {
		return users.User{}, err
	}

	return user, nil
}

func (s *UserService) SignInUser(ctx context.Context, email, password string) (string, error) {
	hashedPassword, err := s.hasher.Hash(password)

	if err != nil {
		return "", err
	}

	userDTO := users.UserSignInDTO{
		Email:    email,
		Password: hashedPassword,
	}

	user, err := s.storage.GetByCredentials(ctx, &userDTO)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
	})

	ss, err := token.SignedString(s.jwtSecret)

	return ss, err
}

func (u UserService) ParseToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return u.jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return "", err
	}

	if !token.Valid || claims.Subject == "" {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}
