package services

import (
	"context"

	"github.com/f4ke-n0name/avito/internal/domain/entities"
	"github.com/f4ke-n0name/avito/internal/domain/errors"
	"github.com/f4ke-n0name/avito/internal/domain/repositories"
	"github.com/f4ke-n0name/avito/internal/domain/services/interfaces"
)

type teamService struct {
	teams repositories.TeamRepository
	users repositories.UserRepository
}

func NewTeamService(teams repositories.TeamRepository, users repositories.UserRepository) interfaces.TeamService {
	return &teamService{teams: teams, users: users}
}

func (s *teamService) CreateTeam(ctx context.Context, team *entities.Team) (*entities.Team, error) {
	existing, _ := s.teams.GetByName(ctx, team.TeamName)
	if existing != nil {
		return nil, errors.ErrTeamExists
	}

	if err := s.teams.Create(ctx, team); err != nil {
		return nil, err
	}
	for _, member := range team.Members {
		_ = s.users.CreateOrUpdate(ctx, &member)
	}

	return team, nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*entities.Team, error) {
	team, err := s.teams.GetByName(ctx, teamName)
	if err != nil || team == nil {
		return nil, errors.ErrTeamNotFound
	}
	return team, nil
}
