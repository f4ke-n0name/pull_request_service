package interfaces

import (
	"context"
	"github.com/f4ke-n0name/avito/internal/domain/entities"
)

type PRService interface {
	CreatePR(ctx context.Context, prID, prName, authorID string) (*entities.PullRequest, error)
	ReplaceReviewer(ctx context.Context, prID, oldReviewerID string) (*entities.PullRequest, string, error)
	Merge(ctx context.Context, prID string) (*entities.PullRequest, error)
	ListByReviewer(ctx context.Context, reviewerID string) ([]entities.PullRequest, error)
}
