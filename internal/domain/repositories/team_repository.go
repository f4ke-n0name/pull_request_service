package repositories

import (
	"context"
	"github.com/f4ke-n0name/avito/internal/domain/entities"
)

type TeamRepository interface {
	Create(ctx context.Context, team *entities.Team) error
	GetByName(ctx context.Context, name string) (*entities.Team, error)
}
