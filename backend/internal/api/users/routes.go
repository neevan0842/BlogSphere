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
	user, err := h.service.getUserByUsername(r.Context(), pgtype.Text{String: username, Valid: true})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}
