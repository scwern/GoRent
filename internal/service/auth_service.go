package service

import (
	"GoRent/internal/domain/user"
	"GoRent/internal/repository"
	"GoRent/pkg/jwt"
	"GoRent/pkg/password"
	"context"
	"errors"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, registerReq *user.RegisterRequest) (*user.User, error)
	Login(ctx context.Context, loginReq *user.LoginRequest) (*user.LoginResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwt      jwt.Manager
}

func NewAuthService(userRepo repository.UserRepository, jwt jwt.Manager) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (s *authService) Register(ctx context.Context, registerReq *user.RegisterRequest) (*user.User, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, registerReq.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := password.Hash(registerReq.Password)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		ID:           uuid.New().String(),
		Name:         registerReq.Name,
		Email:        registerReq.Email,
		PasswordHash: hashedPassword,
		Role:         user.ClientRole,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(ctx context.Context, loginReq *user.LoginRequest) (*user.LoginResponse, error) {
	u, err := s.userRepo.GetByEmail(ctx, loginReq.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !password.Verify(u.PasswordHash, loginReq.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		Token: token,
		User:  u,
	}, nil
}
