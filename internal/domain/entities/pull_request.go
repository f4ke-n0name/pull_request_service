package entities

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PRID      string     `db:"pr_id"`
	Name      string     `db:"pr_name"`
	AuthorID  string     `db:"author_id"`
	Status    PRStatus   `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
	Reviewers []string
}
