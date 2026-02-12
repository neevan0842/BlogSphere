package users

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/internal/common"
	"github.com/neevan0842/BlogSphere/backend/utils"
	"go.uber.org/zap"
)

type handler struct {
	service Service
	logger  *zap.SugaredLogger
	repo    *sqlc.Queries
}

func NewHandler(service Service, logger *zap.SugaredLogger, repo *sqlc.Queries) *handler {
	return &handler{
		service: service,
		logger:  logger,
		repo:    repo,
	}
}

// getRequestingUserID extracts and validates the user ID from the request token
func (h *handler) getRequestingUserID(w http.ResponseWriter, r *http.Request) *pgtype.UUID {
	return common.GetRequestingUserID(r.Context(), w, r, h.repo)
}

func (h *handler) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by authentication middleware)
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("access token in header invalid"))
		return
	}

	userIDUUID, err := utils.StrToUUID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID in token: %s", err.Error()))
		return
	}

	// Fetch user details from the database
	user, err := h.service.getUserByID(r.Context(), userIDUUID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *handler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	userIDUUID, err := utils.StrToUUID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID in token: %s", err.Error()))
		return
	}

	// Fetch user details from the database
	user, err := h.service.getUserByID(r.Context(), userIDUUID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *handler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	// Fetch user details from the database
	user, err := h.service.getUserByUsername(r.Context(), pgtype.Text{String: username, Valid: username != ""})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *handler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var payload UpdateUserRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload: %s", err.Error()))
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload: %s", err.Error()))
		return
	}

	// Get user ID from context (set by authentication middleware)
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("access token in header invalid"))
		return
	}

	// verify user is updating their own profile
	requestedUserID := chi.URLParam(r, "userID")
	if userID != requestedUserID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you can only update your own profile"))
		return
	}

	userIDUUID, err := utils.StrToUUID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID in token: %s", err.Error()))
		return
	}

	// Update user description in the database
	updatedUser, err := h.service.updateUserDescription(r.Context(), userIDUUID, pgtype.Text{String: payload.Description, Valid: payload.Description != ""})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update user description: %s", err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedUser)
}

func (h *handler) HandleDeleteCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by authentication middleware)
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("access token in header invalid"))
		return
	}

	// verify user is deleting their own account
	requestedUserID := chi.URLParam(r, "userID")
	if userID != requestedUserID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you can only delete your own account"))
		return
	}

	userIDUUID, err := utils.StrToUUID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID in token: %s", err.Error()))
		return
	}

	// Delete user from the database
	if err := h.service.deleteUserByID(r.Context(), userIDUUID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete user: %s", err.Error()))
		return
	}

	// Return 204 No Content without response body
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) HandleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	// verify user exists
	_, err := h.service.getUserByUsername(r.Context(), pgtype.Text{String: username, Valid: username != ""})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}

	// check if requester is authenticated
	requestingUserID := h.getRequestingUserID(w, r)

	posts, err := h.service.getPostsByUsername(r.Context(), pgtype.Text{String: username, Valid: username != ""}, requestingUserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch user posts: %s", err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusOK, posts)
}

func (h *handler) HandleGetLikedPosts(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	// verify user exists
	_, err := h.service.getUserByUsername(r.Context(), pgtype.Text{String: username, Valid: username != ""})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}

	// check if requester is authenticated
	requestingUserID := h.getRequestingUserID(w, r)

	posts, err := h.service.getLikedPostsByUsername(r.Context(), pgtype.Text{String: username, Valid: username != ""}, requestingUserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch liked posts: %s", err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusOK, posts)
}
