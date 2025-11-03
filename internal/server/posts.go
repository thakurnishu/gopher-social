package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/thakurnishu/gopher-social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=2000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=2000"`
}

type CreateCommentOnPostPayload struct {
	UserID   int64  `json:"user_id" validate:"required"`
	Content  string `json:"content" validate:"required,max=500"`
	UserName string `json:"username" validate:"required"`
}

type postKey string

const postCtx postKey = "post"

func (s *Server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := r.Context()
	if err := s.store.Posts.Create(ctx, post); err != nil {
		s.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func (s *Server) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := s.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func (s *Server) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := s.store.Posts.Update(r.Context(), post); err != nil {
		s.internalServerError(w, r, err)
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func (s *Server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := getIntParam(r, "postID")
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := s.store.Posts.Delete(ctx, postID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			s.notFoundResponse(w, r, err)
		default:
			s.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) createCommentOnPostHanlder(w http.ResponseWriter, r *http.Request) {
	postID, err := getIntParam(r, "postID")
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	var payload CreateCommentOnPostPayload
	if err := readJSON(w, r, &payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	comment := store.Comment{
		Content:  payload.Content,
		PostID:   postID,
		UserID:   payload.UserID,
		UserName: payload.UserName,
	}

	if err := s.store.Comments.Create(ctx, &comment); err != nil {
		s.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, comment); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func (s *Server) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := getIntParam(r, "postID")
		if err != nil {
			s.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		post, err := s.store.Posts.GetByID(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				s.notFoundResponse(w, r, err)
			default:
				s.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
