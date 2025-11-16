package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thakurnishu/gopher-social/internal/store"
)

type CreateUserPayload struct {
	Username  string `json:"username" validate:"required,max=50"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"-"`
}

type UpdateUserPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=2000"`
}


//func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
//	var payload CreateUserPayload
//	if err := readJSON(w, r, &payload); err != nil {
//		s.badRequestResponse(w, r, err)
//		return
//	}
//
//	if err := Validate.Struct(payload); err != nil {
//		s.badRequestResponse(w, r, err)
//		return
//	}
//
//	post := &store.User{
//		Username: "",
//	}
//
//	ctx := r.Context()
//	if err := s.store.Posts.Create(ctx, post); err != nil {
//		s.internalServerError(w, r, err)
//		return
//	}
//
//	if err := s.jsonResponse(w, http.StatusCreated, post); err != nil {
//		s.internalServerError(w, r, err)
//		return
//	}
//}

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := s.store.Users.GetByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			s.notFoundResponse(w, r, err)
		default:
			s.internalServerError(w, r, err)
		}
		return
	}


	if err := s.jsonResponse(w, http.StatusOK, user); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

//func (s *Server) updatePostHandler(w http.ResponseWriter, r *http.Request) {
//	post := getPostFromCtx(r)
//
//	var payload UpdatePostPayload
//	if err := readJSON(w, r, &payload); err != nil {
//		s.badRequestResponse(w, r, err)
//		return
//	}
//	if err := Validate.Struct(payload); err != nil {
//		s.badRequestResponse(w, r, err)
//		return
//	}
//
//	if payload.Title != nil {
//		post.Title = *payload.Title
//	}
//	if payload.Content != nil {
//		post.Content = *payload.Content
//	}
//
//	if err := s.store.Posts.Update(r.Context(), post); err != nil {
//		s.internalServerError(w, r, err)
//		return
//	}
//
//	if err := s.jsonResponse(w, http.StatusOK, post); err != nil {
//		s.internalServerError(w, r, err)
//		return
//	}
//}
//
//func (s *Server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
//	postID, err := getIntParam(r, "postID")
//	if err != nil {
//		s.internalServerError(w, r, err)
//		return
//	}
//
//	ctx := r.Context()
//	if err := s.store.Posts.Delete(ctx, postID); err != nil {
//		switch {
//		case errors.Is(err, store.ErrNotFound):
//			s.notFoundResponse(w, r, err)
//		default:
//			s.internalServerError(w, r, err)
//		}
//		return
//	}
//
//	w.WriteHeader(http.StatusNoContent)
//}
//
//func (s *Server) createCommentOnPostHanlder(w http.ResponseWriter, r *http.Request) {
//	postID, err := getIntParam(r, "postID")
//	if err != nil {
//		s.internalServerError(w, r, err)
//		return
//	}
//
//	var payload CreateCommentOnPostPayload
//	if err := readJSON(w, r, &payload); err != nil {
//		s.badRequestResponse(w, r, err)
//		return
//	}
//
//	if err := Validate.Struct(payload); err != nil {
//		s.badRequestResponse(w, r, err)
//		return
//	}
//
//	ctx := r.Context()
//	comment := store.Comment{
//		Content:  payload.Content,
//		PostID:   postID,
//		UserID:   payload.UserID,
//		UserName: payload.UserName,
//	}
//
//	if err := s.store.Comments.Create(ctx, &comment); err != nil {
//		s.internalServerError(w, r, err)
//		return
//	}
//
//	if err := s.jsonResponse(w, http.StatusCreated, comment); err != nil {
//		s.internalServerError(w, r, err)
//		return
//	}
//}
