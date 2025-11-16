package services

import (
	"context"
	"github.com/f4ke-n0name/avito/internal/domain/entities"
	"github.com/f4ke-n0name/avito/internal/domain/errors"
	"github.com/f4ke-n0name/avito/internal/domain/repositories"
	"github.com/f4ke-n0name/avito/internal/domain/services/interfaces"
	"github.com/jackc/pgconn"
	"math/rand"
	"time"
)

type prService struct {
	users repositories.UserRepository
	teams repositories.TeamRepository
	prs   repositories.PullRequestRepository

	withTx func(ctx context.Context, fn func(txCtx context.Context) error) error
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewPRService(
	users repositories.UserRepository,
	teams repositories.TeamRepository,
	prs repositories.PullRequestRepository,
	withTx func(ctx context.Context, fn func(txCtx context.Context) error) error,
) interfaces.PRService {
	return &prService{
		users:  users,
		teams:  teams,
		prs:    prs,
		withTx: withTx,
	}
}

func (s *prService) CreatePR(ctx context.Context, prID, prName, authorID string) (*entities.PullRequest, error) {
	author, err := s.users.GetByID(ctx, authorID)
	if err != nil || author == nil {
		return nil, errors.ErrUserNotFound
	}

	var pr *entities.PullRequest

	err = s.withTx(ctx, func(txCtx context.Context) error {
		pr = &entities.PullRequest{
			PRID:     prID,
			Name:     prName,
			AuthorID: authorID,
			Status:   entities.PRStatusOpen,
		}

		if err := s.prs.Create(txCtx, pr); err != nil {
			if IsUniqueViolation(err) {
				return errors.ErrPRExists
			}
			return err
		}

		candidates, err := s.users.ListActiveByTeam(txCtx, author.TeamName)
		if err != nil {
			return err
		}

		var filtered []entities.User
		for _, u := range candidates {
			if u.UserID != authorID {
				filtered = append(filtered, u)
			}
		}

		reviewers := pickTwoRandom(filtered)

		var ids []string
		for _, r := range reviewers {
			ids = append(ids, r.UserID)
		}

		if len(ids) > 0 {
			if err := s.prs.AssignReviewers(txCtx, prID, ids); err != nil {
				return err
			}
			pr.Reviewers = ids
		}

		return nil
	})

	return pr, err
}

func (s *prService) ReplaceReviewer(ctx context.Context, prID, oldReviewerID string) (*entities.PullRequest, string, error) {
	pr, err := s.prs.GetByID(ctx, prID)
	if err != nil || pr == nil {
		return nil, "", errors.ErrPRNotFound
	}
	if pr.Status == entities.PRStatusMerged {
		return nil, "", errors.ErrPRAlreadyMerged
	}
	assigned := false
	for _, r := range pr.Reviewers {
		if r == oldReviewerID {
			assigned = true
			break
		}
	}
	if !assigned {
		return nil, "", errors.ErrNoSuchReviewer
	}
	oldReviewer, err := s.users.GetByID(ctx, oldReviewerID)
	if err != nil || oldReviewer == nil {
		return nil, "", errors.ErrUserNotFound
	}
	var updated *entities.PullRequest
	var newID string
	err = s.withTx(ctx, func(txCtx context.Context) error {
		candidates, err := s.users.ListActiveByTeam(txCtx, oldReviewer.TeamName)
		if err != nil {
			return err
		}
		var filtered []entities.User
		for _, u := range candidates {
			if u.UserID != oldReviewerID {
				filtered = append(filtered, u)
			}
		}
		if len(filtered) == 0 {
			return errors.ErrNoCandidates
		}
		newUser := filtered[rnd.Intn(len(filtered))]
		newID = newUser.UserID
		if err := s.prs.ReplaceReviewer(txCtx, prID, oldReviewerID, newID); err != nil {
			return err
		}
		updated, err = s.prs.GetByID(txCtx, prID)
		return err
	})
	return updated, newID, err
}

func (s *prService) Merge(ctx context.Context, prID string) (*entities.PullRequest, error) {
	pr, err := s.prs.GetByID(ctx, prID)
	if err != nil || pr == nil {
		return nil, errors.ErrPRNotFound
	}
	if pr.Status == entities.PRStatusMerged {
		return pr, nil
	}
	if err := s.prs.MarkMerged(ctx, prID); err != nil {
		return nil, err
	}
	return s.prs.GetByID(ctx, prID)
}

func (s *prService) ListByReviewer(ctx context.Context, reviewerID string) ([]entities.PullRequest, error) {
	return s.prs.ListByReviewer(ctx, reviewerID)
}

func pickTwoRandom(users []entities.User) []entities.User {
	n := len(users)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []entities.User{users[0]}
	}
	i := rnd.Intn(n)
	j := rnd.Intn(n - 1)
	if j >= i {
		j++
	}
	return []entities.User{users[i], users[j]}
}

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505"
	}
	return false
}
