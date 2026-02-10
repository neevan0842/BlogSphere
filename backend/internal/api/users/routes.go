package users

import (
	"fmt"
	"net/http"

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

func (h *handler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
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
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get user: %s", err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}
