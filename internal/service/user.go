package service

import (
	"context"

	"github.com/gin-gonic/gin"
)

type UserStorage interface {
	Save(ctx context.Context, name, email, password string) (string, error)
	Get(ctx context.Context, id int) (string, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserService struct {
  storage UserStorage
	hasher  PasswordHasher
}

func NewUserService(storage UserStorage, hasher PasswordHasher) *UserService {
	return &UserService{storage: storage, hasher: hasher}
}

func (s *UserService) SignUp(ctx *gin.Context, name, email, password string) (string, error) {
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return "", err
	}
	return s.storage.Save(ctx, name, email, hashedPassword)
}

func (s *UserService) GetUser(ctx *gin.Context, id int) (string, error) {
	return s.storage.Get(ctx, id)
}