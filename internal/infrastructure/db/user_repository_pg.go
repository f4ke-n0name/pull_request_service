package db

import (
	"context"

	"github.com/f4ke-n0name/avito/internal/domain/entities"
	"github.com/f4ke-n0name/avito/internal/domain/repositories"
	"github.com/jackc/pgx/v5"
)

type UserRepositoryPG struct {
	db *PG
}

func NewUserRepositoryPG(pg *PG) repositories.UserRepository {
	return &UserRepositoryPG{db: pg}
}

func (r *UserRepositoryPG) querier(ctx context.Context) dbQuerier {
	if tx, ok := TxFromContext(ctx); ok && tx != nil {
		return tx
	}
	return r.db.Pool
}

func (r *UserRepositoryPG) CreateOrUpdate(ctx context.Context, u *entities.User) error {
	q := `
        INSERT INTO users (user_id, username, is_active, team_name)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE
        SET username = EXCLUDED.username,
            is_active = EXCLUDED.is_active,
            team_name = EXCLUDED.team_name
    `
	_, err := r.querier(ctx).Exec(ctx, q, u.UserID, u.Username, u.IsActive, u.TeamName)
	return err
}

func (r *UserRepositoryPG) GetByID(ctx context.Context, id string) (*entities.User, error) {
	q := `SELECT user_id, username, is_active, team_name FROM users WHERE user_id = $1`

	u := &entities.User{}
	err := r.querier(ctx).QueryRow(ctx, q, id).Scan(&u.UserID, &u.Username, &u.IsActive, &u.TeamName)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepositoryPG) ListByTeam(ctx context.Context, team string) ([]entities.User, error) {
	q := `
        SELECT user_id, username, is_active, team_name
        FROM users
        WHERE team_name = $1
    `

	rows, err := r.querier(ctx).Query(ctx, q, team)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.User
	for rows.Next() {
		var u entities.User
		_ = rows.Scan(&u.UserID, &u.Username, &u.IsActive, &u.TeamName)
		result = append(result, u)
	}
	return result, nil
}

func (r *UserRepositoryPG) ListActiveByTeam(ctx context.Context, team string) ([]entities.User, error) {
	q := `
        SELECT user_id, username, is_active, team_name
        FROM users
        WHERE team_name = $1 AND is_active = true
    `

	rows, err := r.querier(ctx).Query(ctx, q, team)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []entities.User
	for rows.Next() {
		var u entities.User
		_ = rows.Scan(&u.UserID, &u.Username, &u.IsActive, &u.TeamName)
		res = append(res, u)
	}
	return res, nil
}

func (r *UserRepositoryPG) SetActive(ctx context.Context, id string, active bool) error {
	q := `UPDATE users SET is_active = $1 WHERE user_id = $2`
	_, err := r.querier(ctx).Exec(ctx, q, active, id)
	return err
}
