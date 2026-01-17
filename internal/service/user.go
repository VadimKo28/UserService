package service

import (
	"context"
	"user_advt/internal/domain/users"

	"github.com/gin-gonic/gin"
)

type UserStorage interface {
	Save(ctx context.Context, user *users.UserCreateDTO) (int, error)
	Get(ctx context.Context, id string) (users.GetUserDTO, error)
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

func (s *UserService) SignUp(ctx *gin.Context, name, email, password string) (int, error) {
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

func (s *UserService) GetUser(ctx *gin.Context, id string) (users.GetUserDTO, error) {

	user, err := s.storage.Get(ctx, id)

	if err != nil {
		return users.GetUserDTO{}, err
	}

	return user, nil
}