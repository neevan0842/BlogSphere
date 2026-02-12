package comments

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (h *handler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "commentID")
	commentIDUUID, _ := utils.StrToUUID(commentID)

	// Get user ID from context (set by authentication middleware)
	userID, _ := utils.GetUserIDFromContext(r.Context())

	comment, err := h.repo.GetCommentByID(r.Context(), commentIDUUID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Only allow comment author to delete the comment
	if comment.UserID.String() != userID {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}

	if err := h.service.HandleDeleteComment(r.Context(), commentID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) HandleUpdateComment(w http.ResponseWriter, r *http.Request) {}
