package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thakurnishu/gopher-social/internal/store"
)

type userKey string

const userCtx userKey = "user"

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
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
	user, err := getUserFromContext(r)
	if err != nil {
		// client likely hit an endpoint not wrapped by middleware or we have a server bug
		s.internalServerError(w, r, err)
		return
	}

	if err := s.jsonResponse(w, http.StatusOK, user); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (s *Server) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser, err := getUserFromContext(r)
	if err != nil {
		// client likely hit an endpoint not wrapped by middleware or we have a server bug
		s.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	// revert back to auth userid from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	err = s.store.Followers.Follow(ctx, payload.UserID, followedUser.ID)
	if err != nil {
		switch err {
		case store.ErrAlreadyFollowing:
			s.confictResponse(w, r, err)
		default:
			s.internalServerError(w, r, err)
		}
		return
	}

	if err := s.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func (s *Server) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser, err := getUserFromContext(r)
	if err != nil {
		// client likely hit an endpoint not wrapped by middleware or we have a server bug
		s.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	// revert back to auth userid from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	err = s.store.Followers.UnFollow(ctx, payload.UserID, followedUser.ID)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	if err := s.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		s.internalServerError(w, r, err)
		return
	}
}

func (s *Server) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				return
			default:
				s.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) (*store.User, error) {
	v := r.Context().Value(userCtx)
	if v == nil {
		return nil, fmt.Errorf("user not found in context")
	}

	u, ok := v.(*store.User)
	if !ok || u == nil {
		return nil, fmt.Errorf("user in context has wrong type or is nil")
	}

	return u, nil
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
