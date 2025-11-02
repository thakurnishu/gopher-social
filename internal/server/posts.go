package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thakurnishu/gopher-social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (s *Server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		s.badRequestResponse(w, r, err)
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
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
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	post, err := s.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			s.notFoundResponse(w, r, err)
		default:
			s.internalServerError(w, r, err)
		}
		return
	}

	comments, err := s.store.Comments.GetByPostID(ctx, id)
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
