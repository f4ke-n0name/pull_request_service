package db

import (
	"context"

	"github.com/f4ke-n0name/avito/internal/domain/entities"
	"github.com/f4ke-n0name/avito/internal/domain/repositories"
	"github.com/jackc/pgx/v5"
)

type TeamRepositoryPG struct {
	db *PG
}

func NewTeamRepositoryPG(db *PG) repositories.TeamRepository {
	return &TeamRepositoryPG{db: db}
}

func (r *TeamRepositoryPG) querier(ctx context.Context) dbQuerier {
	if tx, ok := TxFromContext(ctx); ok && tx != nil {
		return tx
	}
	return r.db.Pool
}

func (r *TeamRepositoryPG) Create(ctx context.Context, t *entities.Team) error {
	qTeam := `INSERT INTO teams (team_name) VALUES ($1)`
	if _, err := r.querier(ctx).Exec(ctx, qTeam, t.TeamName); err != nil {
		return err
	}
	qUser := `INSERT INTO users (user_id, username, is_active, team_name) VALUES ($1, $2, $3, $4)
              ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username, is_active = EXCLUDED.is_active, team_name = EXCLUDED.team_name`
	for _, member := range t.Members {
		if _, err := r.querier(ctx).Exec(ctx, qUser, member.UserID, member.Username, member.IsActive, t.TeamName); err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepositoryPG) GetByName(ctx context.Context, name string) (*entities.Team, error) {
	var tmp string
	err := r.querier(ctx).QueryRow(ctx, `SELECT team_name FROM teams WHERE team_name=$1`, name).Scan(&tmp)
	if err != nil {
		if pgx.ErrNoRows == err {
			return nil, nil
		}
		return nil, err
	}
	q := `SELECT user_id, username, is_active FROM users WHERE team_name=$1`
	rows, err := r.querier(ctx).Query(ctx, q, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	team := &entities.Team{TeamName: name}
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.IsActive); err != nil {
			return nil, err
		}
		user.TeamName = name
		team.Members = append(team.Members, user)
	}

	return team, nil
}
