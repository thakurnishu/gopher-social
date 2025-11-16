package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrAlreadyFollowing = errors.New("user already follows this user")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, postID int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, postID int64) error
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, userID int64) (*User, error)
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerId int64, userID int64) error
		UnFollow(ctx context.Context, followerId int64, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
