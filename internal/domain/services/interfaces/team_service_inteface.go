package interfaces

import (
	"context"

	"github.com/f4ke-n0name/avito/internal/domain/entities"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *entities.Team) (*entities.Team, error)
	GetTeam(ctx context.Context, teamName string) (*entities.Team, error)
}
