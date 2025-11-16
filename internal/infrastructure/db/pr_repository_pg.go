package db

import (
	"context"
	"github.com/f4ke-n0name/avito/internal/domain/entities"
	"github.com/f4ke-n0name/avito/internal/domain/repositories"
	"github.com/jackc/pgx/v5"
)

type PRRepositoryPG struct {
	db *PG
}

func NewPRRepositoryPG(db *PG) repositories.PullRequestRepository {
	return &PRRepositoryPG{db: db}
}

func (r *PRRepositoryPG) querier(ctx context.Context) dbQuerier {
	if tx, ok := TxFromContext(ctx); ok && tx != nil {
		return tx
	}
	return r.db.Pool
}

func (r *PRRepositoryPG) Create(ctx context.Context, pr *entities.PullRequest) error {
	q := `
        INSERT INTO pull_requests (pr_id, pr_name, author_id, status)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.querier(ctx).Exec(ctx, q, pr.PRID, pr.Name, pr.AuthorID, pr.Status)
	return err
}

func (r *PRRepositoryPG) GetByID(ctx context.Context, id string) (*entities.PullRequest, error) {
	pr := &entities.PullRequest{}
	q := `
        SELECT pr_id, pr_name, author_id, status, created_at, merged_at
        FROM pull_requests
        WHERE pr_id = $1
    `
	err := r.querier(ctx).QueryRow(ctx, q, id).Scan(&pr.PRID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	q2 := `
        SELECT reviewer_id
        FROM pull_request_reviewers
        WHERE pr_id = $1
        ORDER BY assigned_at
    `
	rows, err := r.querier(ctx).Query(ctx, q2, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reviewer string
		_ = rows.Scan(&reviewer)
		pr.Reviewers = append(pr.Reviewers, reviewer)
	}
	return pr, nil
}

func (r *PRRepositoryPG) ListByReviewer(ctx context.Context, reviewerID string) ([]entities.PullRequest, error) {
	q := `
        SELECT pr.pr_id, pr.pr_name, pr.author_id, pr.status, pr.created_at, pr.merged_at,
               ARRAY_AGG(rr.reviewer_id) AS reviewers
        FROM pull_requests pr
        JOIN pull_request_reviewers rr ON rr.pr_id = pr.pr_id
        WHERE pr.pr_id IN (
            SELECT pr_id FROM pull_request_reviewers WHERE reviewer_id = $1
        )
        GROUP BY pr.pr_id
        ORDER BY pr.created_at DESC
    `
	rows, err := r.querier(ctx).Query(ctx, q, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.PullRequest
	for rows.Next() {
		var pr entities.PullRequest
		var reviewers []string
		if err := rows.Scan(&pr.PRID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt, &reviewers); err != nil {
			return nil, err
		}
		pr.Reviewers = reviewers
		result = append(result, pr)
	}

	return result, nil
}

func (r *PRRepositoryPG) AssignReviewers(ctx context.Context, prID string, reviewers []string) error {
	q := `
        INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
    `
	for _, reviewerID := range reviewers {
		if _, err := r.querier(ctx).Exec(ctx, q, prID, reviewerID); err != nil {
			return err
		}
	}
	return nil
}

func (r *PRRepositoryPG) ReplaceReviewer(ctx context.Context, prID string, oldID, newID string) error {
	if _, err := r.querier(ctx).Exec(ctx,
		`DELETE FROM pull_request_reviewers WHERE pr_id = $1 AND reviewer_id = $2`,
		prID, oldID); err != nil {
		return err
	}
	if _, err := r.querier(ctx).Exec(ctx,
		`INSERT INTO pull_request_reviewers (pr_id, reviewer_id) VALUES ($1, $2)`,
		prID, newID); err != nil {
		return err
	}
	return nil
}

func (r *PRRepositoryPG) MarkMerged(ctx context.Context, prID string) error {
	_, err := r.querier(ctx).Exec(ctx,
		`UPDATE pull_requests SET status='MERGED', merged_at=now() WHERE pr_id=$1`,
		prID)
	return err
}
