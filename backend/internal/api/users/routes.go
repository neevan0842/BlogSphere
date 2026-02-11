package users

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/utils"
	"go.uber.org/zap"
)

type handler struct {
	service Service
	logger  *zap.SugaredLogger
}

func NewHandler(service Service, logger *zap.SugaredLogger) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
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
	userID := chi.URLParam(r, "id")

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

func (h *handler) HandleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	// verify user exists
	_, err := h.service.getUserByUsername(r.Context(), pgtype.Text{String: username, Valid: username != ""})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}

	posts, err := h.service.getPostsByUsername(r.Context(), pgtype.Text{String: username, Valid: username != ""})
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

	// check if requester is authenticated to determine whether to include "liked_by_requester" field
	var requestingUserID *pgtype.UUID
	token := utils.GetTokenFromRequest(r)
	userIDUUID, err := utils.GetUserIDFromToken(w, token)
	if err != nil {
		requestingUserID = nil // requester is not authenticated
	} else {
		user, err := h.service.getUserByID(r.Context(), userIDUUID)
		if err != nil {
			requestingUserID = nil // requester is not authenticated
		} else {
			requestingUserID = &user.ID
		}
	}

	posts, err := h.service.getLikedPostsByUsername(r.Context(), pgtype.Text{String: username, Valid: username != ""}, requestingUserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch liked posts: %s", err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusOK, posts)
}
