package repositories

import (
	"context"
	"github.com/f4ke-n0name/avito/internal/domain/entities"
)

type UserRepository interface {
	CreateOrUpdate(ctx context.Context, u *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	ListByTeam(ctx context.Context, team string) ([]entities.User, error)
	ListActiveByTeam(ctx context.Context, team string) ([]entities.User, error)
	SetActive(ctx context.Context, id string, active bool) error
}
