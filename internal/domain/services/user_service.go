package services

import (
	"context"
	"github.com/f4ke-n0name/avito/internal/domain/entities"
	"github.com/f4ke-n0name/avito/internal/domain/repositories"
	"github.com/f4ke-n0name/avito/internal/domain/services/interfaces"
	"github.com/f4ke-n0name/avito/internal/domain/errors"
)

type userService struct {
	users repositories.UserRepository
}

func NewUserService(users repositories.UserRepository) interfaces.UserService {
	return &userService{users: users}
}

func (s *userService) SetIsActive(ctx context.Context, userID string, isActive bool) (*entities.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.ErrUserNotFound
	}

	if err := s.users.SetActive(ctx, userID, isActive); err != nil {
		return nil, err
	}

	user.IsActive = isActive
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, userID string) (*entities.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) ListByTeam(ctx context.Context, teamName string) ([]entities.User, error) {
	return s.users.ListByTeam(ctx, teamName)
}

func (s *userService) ListActiveByTeam(ctx context.Context, teamName string) ([]entities.User, error) {
	return s.users.ListActiveByTeam(ctx, teamName)
}
