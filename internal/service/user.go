package service

import (
	"context"
)

type UserStorage interface {
	SaveUser(ctx context.Context, name, email, password string) (string, error)
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

func (s *UserService) SignUp(ctx context.Context, name, email, password string) (string, error) {
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return "", err
	}
	return s.storage.SaveUser(ctx, name, email, hashedPassword)
}