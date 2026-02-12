package comments

import (
	"net/http"

	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
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

func (h *handler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Get requesting user ID (set by authentication middleware)
	userID, _ := utils.GetUserIDFromContext(r.Context())

	comment, err := h.service.HandleCreateComment(r.Context(), payload.PostID, userID, payload.Body)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, comment)
}

func (h *handler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {}

func (h *handler) HandleUpdateComment(w http.ResponseWriter, r *http.Request) {}
