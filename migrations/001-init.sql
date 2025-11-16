BEGIN;

CREATE TABLE teams (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    team_name TEXT NOT NULL,
    CONSTRAINT fk_users_team FOREIGN KEY (team_name)
        REFERENCES teams (team_name)
        ON DELETE RESTRICT
);

CREATE INDEX idx_users_team_name ON users(team_name);

CREATE TABLE pull_requests (
    pr_id TEXT PRIMARY KEY,
    pr_name TEXT NOT NULL,
    author_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    merged_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_pr_author FOREIGN KEY (author_id)
        REFERENCES users (user_id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_pr_author_id ON pull_requests(author_id);
CREATE INDEX idx_pr_status ON pull_requests(status);

CREATE TABLE pull_request_reviewers (
    pr_id TEXT NOT NULL,
    reviewer_id TEXT NOT NULL,
    assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (pr_id, reviewer_id),
    CONSTRAINT fk_prrev_pr FOREIGN KEY (pr_id)
        REFERENCES pull_requests (pr_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_prrev_user FOREIGN KEY (reviewer_id)
        REFERENCES users (user_id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_prrev_reviewer_id ON pull_request_reviewers(reviewer_id);
CREATE INDEX idx_prrev_pr_id ON pull_request_reviewers(pr_id);

COMMIT;
