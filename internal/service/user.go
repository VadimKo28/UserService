package service

import (
	"context"
	"user_advt/internal/domain/users"
)

type UserStorage interface {
	Save(ctx context.Context, user *users.UserCreateDTO) (int, error)
	Get(ctx context.Context, id string) (users.User, error)
	GetByCredentials(ctx context.Context, user *users.UserSignInDTO) (users.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type TokenSigner interface {
	SignToken(userID int) (string, error)
}

type UserService struct {
	storage      UserStorage
	hasher       PasswordHasher
	tokenService TokenSigner
}

func NewUserService(storage UserStorage, hasher PasswordHasher, tokenService TokenSigner) *UserService {
	return &UserService{storage: storage, hasher: hasher, tokenService: tokenService}
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

	return s.tokenService.SignToken(int(user.ID))
}
