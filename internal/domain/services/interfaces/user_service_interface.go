package interfaces

import (
	"context"
	"github.com/f4ke-n0name/avito/internal/domain/entities"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*entities.User, error)
	GetByID(ctx context.Context, userID string) (*entities.User, error)
	ListByTeam(ctx context.Context, teamName string) ([]entities.User, error)
	ListActiveByTeam(ctx context.Context, teamName string) ([]entities.User, error)
}
