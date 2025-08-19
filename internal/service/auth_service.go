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
	Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error)
	Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error)
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

func (s *authService) Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         user.ClientRole,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	u, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !password.Verify(u.PasswordHash, req.Password) {
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
