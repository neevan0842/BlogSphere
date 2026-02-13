package categories

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

func (h *handler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, categories)
}
