package service

import (
	"app/internal/domain/users"
	"app/internal/storage"
	"context"
	"time"
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
	GenerateTokens(userID int) (string, string, error)
}

type RefreshTokenStorage interface {
	Save(ctx context.Context, userID int, token string, expiresAt time.Time) error
	Get(ctx context.Context, token string) (int, time.Time, error)
	Delete(ctx context.Context, token string) error
}

type UserService struct {
	storage         UserStorage
	refreshStorage  RefreshTokenStorage
	hasher          PasswordHasher
	tokenService    TokenSigner
	refreshTokenTTL time.Duration
}

func NewUserService(
	storage UserStorage,
	refreshStorage RefreshTokenStorage,
	hasher PasswordHasher,
	tokenService TokenSigner,
	refreshTokenTTL time.Duration,
) *UserService {
	return &UserService{
		storage:         storage,
		refreshStorage:  refreshStorage,
		hasher:          hasher,
		tokenService:    tokenService,
		refreshTokenTTL: refreshTokenTTL,
	}
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

func (s *UserService) SignInUser(ctx context.Context, email, password string) (string, string, error) {
	hashedPassword, err := s.hasher.Hash(password)

	if err != nil {
		return "", "", err
	}

	userDTO := users.UserSignInDTO{
		Email:    email,
		Password: hashedPassword,
	}

	user, err := s.storage.GetByCredentials(ctx, &userDTO)

	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := s.tokenService.GenerateTokens(int(user.ID))

	if err != nil {
		return "", "", err
	}

	if err := s.refreshStorage.Save(ctx, int(user.ID), refreshToken, time.Now().Add(s.refreshTokenTTL)); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	userID, expiresAt, err := s.refreshStorage.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if time.Now().After(expiresAt) {
		_ = s.refreshStorage.Delete(ctx, refreshToken)
		return "", "", storage.ErrRefreshTokenExpired
	}

	if err := s.refreshStorage.Delete(ctx, refreshToken); err != nil {
		return "", "", err
	}

	accessToken, newRefreshToken, err := s.tokenService.GenerateTokens(userID)
	if err != nil {
		return "", "", err
	}

	if err := s.refreshStorage.Save(ctx, userID, newRefreshToken, time.Now().Add(s.refreshTokenTTL)); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *UserService) LogOut(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return storage.ErrRefreshTokenNotFound
	}

	return s.refreshStorage.Delete(ctx, refreshToken)
}
