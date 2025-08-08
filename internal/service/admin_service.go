package service

import (
	"GoRent/internal/domain/user"
	"GoRent/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
)

type AdminService interface {
	ChangeUserRole(ctx context.Context, userID string, req *user.ChangeRoleRequest) error
	ListUsers(ctx context.Context) ([]*user.User, error)
}

type adminService struct {
	userRepo repository.UserRepository
}

func NewAdminService(userRepo repository.UserRepository) AdminService {
	return &adminService{
		userRepo: userRepo,
	}
}

func (s *adminService) ChangeUserRole(ctx context.Context, userID string, req *user.ChangeRoleRequest) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.UpdateRole(ctx, u.ID, req.Role)
}

func (s *adminService) ListUsers(ctx context.Context) ([]*user.User, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}
