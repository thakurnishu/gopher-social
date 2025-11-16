package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	FollowerID int64  `json:"follower_id"`
	UserID     int64  `json:"user_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerId int64, userID int64) error {
	query := `
		INSERT INTO followers (follower_id, user_id)
		VALUES ($1, $2)
	`

	ctx, canel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer canel()

	_, err := s.db.ExecContext(ctx, query, followerId, userID)
	if err != nil {
		// Check for PostgreSQL unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // unique_violation
			return ErrAlreadyFollowing
		}
		return err
	}

	return nil
}

func (s *FollowerStore) UnFollow(ctx context.Context, followerId int64, userID int64) error {
	query := `
		DELETE FROM followers 
		WHERE follower_id = $1 AND user_id = $2
	`

	ctx, canel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer canel()

	_, err := s.db.ExecContext(ctx, query, followerId, userID)
	return err
}
