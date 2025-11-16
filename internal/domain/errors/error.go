package errors

import "errors"

var (
	ErrPRNotFound        = errors.New("pull request not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrTeamNotFound      = errors.New("team not found")
	ErrPRAlreadyMerged   = errors.New("pull request already merged")
	ErrReviewerNotInTeam = errors.New("reviewer is not in expected team")
	ErrReviewerInactive  = errors.New("reviewer is inactive")
	ErrNoCandidates      = errors.New("no active candidates in team")
	ErrNoSuchReviewer    = errors.New("this reviewer is not assigned to PR")
	ErrPRExists          = errors.New("pull request already exists")
	ErrTeamExists        = errors.New("team already exists")
)
