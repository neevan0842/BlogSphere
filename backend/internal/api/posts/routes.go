package posts

import (
	"fmt"
	"net/http"
	"strconv"

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

func (h *handler) HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	searchStr := r.URL.Query().Get("search")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Get requesting user ID (if authenticated)
	requestingUserID := h.getRequestingUserID(w, r)

	// Fetch posts with pagination and search
	posts, err := h.service.getPostsPaginated(r.Context(), searchStr, limit, offset, requestingUserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch posts: %s", err.Error()))
		return
	}

	hasMore := len(posts) == limit
	result := PaginatedResponse{
		Posts:   posts,
		Page:    page,
		Limit:   limit,
		HasMore: hasMore,
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

// getRequestingUserID extracts and validates the user ID from the request token
func (h *handler) getRequestingUserID(w http.ResponseWriter, r *http.Request) *pgtype.UUID {
	return common.GetRequestingUserID(r.Context(), w, r, h.repo)
}

func (h *handler) HandleGetPostsBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	requestingUserID := h.getRequestingUserID(w, r)

	post, err := h.service.getPostBySlug(r.Context(), slug, requestingUserID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("failed to fetch post: %s", err.Error()))
		return
	}
	utils.WriteJSON(w, http.StatusOK, post)
}

func (h *handler) HandleGetCommentsByPostSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	comments, err := h.service.getCommentsByPostSlug(r.Context(), slug)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch comments: %s", err.Error()))
		return
	}
	utils.WriteJSON(w, http.StatusOK, comments)
}
