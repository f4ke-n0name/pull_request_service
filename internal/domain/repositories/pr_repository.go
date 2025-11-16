package repositories

import (
	"context"
	"github.com/f4ke-n0name/avito/internal/domain/entities"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *entities.PullRequest) error
	GetByID(ctx context.Context, id string) (*entities.PullRequest, error)
	ListByReviewer(ctx context.Context, reviewerID string) ([]entities.PullRequest, error)
	AssignReviewers(ctx context.Context, prID string, reviewers []string) error
	ReplaceReviewer(ctx context.Context, prID string, oldID, newID string) error
	MarkMerged(ctx context.Context, prID string) error
}
