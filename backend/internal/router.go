package internal

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	envs "github.com/neevan0842/BlogSphere/backend/config"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/internal/api/auth"
	"github.com/neevan0842/BlogSphere/backend/internal/api/categories"
	"github.com/neevan0842/BlogSphere/backend/internal/api/comments"
	"github.com/neevan0842/BlogSphere/backend/internal/api/posts"
	"github.com/neevan0842/BlogSphere/backend/internal/api/users"

	mw "github.com/neevan0842/BlogSphere/backend/internal/middleware"
	"github.com/neevan0842/BlogSphere/backend/utils"
	"go.uber.org/zap"
)

type application struct {
	config config
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

type config struct {
	addr string
	dsn  string
}

func NewAPIServer(addr string, db *pgxpool.Pool, logger *zap.SugaredLogger) *application {
	return &application{
		config: config{
			addr: addr,
			dsn:  envs.Envs.DATABASE_URL,
		},
		db:     db,
		logger: logger,
	}
}

func (app *application) Run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infof("server has started at http://localhost%s", app.config.addr)

	return srv.ListenAndServe()
}

func (app *application) Mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{envs.Envs.CORS_ALLOWED_ORIGIN},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Initialize services and handlers
	repo := sqlc.New(app.db)

	authService := auth.NewService(repo, app.db)
	authHandler := auth.NewHandler(authService, app.logger)

	userService := users.NewService(repo, app.db)
	userHandler := users.NewHandler(userService, app.logger, repo)

	postService := posts.NewService(repo, app.db)
	postHandler := posts.NewHandler(postService, app.logger, repo)

	commentService := comments.NewService(repo, app.db)
	commentHandler := comments.NewHandler(commentService, app.logger, repo)

	categoryService := categories.NewService(repo, app.db)
	categoryHandler := categories.NewHandler(categoryService, app.logger, repo)

	// Initialize middleware
	authMiddleware := mw.NewMiddleware(repo, app.logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {

		// auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Get("/google", authHandler.HandleGoogleLogin)
			r.Get("/google/callback", authHandler.HandleGoogleAuthCallback)
			r.Post("/refresh", authHandler.HandleRefresh)
		})

		// user routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/u/{username}", userHandler.HandleGetUserByUsername)
			r.Get("/u/{username}/posts", userHandler.HandleGetUserPosts)
			r.Get("/u/{username}/liked-posts", userHandler.HandleGetLikedPosts)
			r.Get("/{userID}", userHandler.HandleGetUserByID)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.UserAuthentication) // Apply authentication middleware to all /users routes
				r.Get("/me", userHandler.HandleGetCurrentUser)
				r.Patch("/{userID}", userHandler.HandleUpdateUser)
				r.Delete("/{userID}", userHandler.HandleDeleteCurrentUser)
			})
		})

		// post routes
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", postHandler.HandleGetPosts)
			r.Get("/{slug}", postHandler.HandleGetPostsBySlug)
			r.Get("/{slug}/comments", postHandler.HandleGetCommentsByPostSlug)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.UserAuthentication)
				r.Post("/", postHandler.HandleCreatePost)
				r.Put("/{postID}", postHandler.HandleUpdatePost)
				r.Delete("/{postID}", postHandler.HandleDeletePost)
				r.Post("/{postID}/likes", postHandler.HandlePostLikes)
			})
		})

		// comment routes
		r.Route("/comments", func(r chi.Router) {
			r.Use(authMiddleware.UserAuthentication)
			r.Post("/", commentHandler.HandleCreateComment)
			r.Delete("/{commentID}", commentHandler.HandleDeleteComment)
			r.Patch("/{commentID}", commentHandler.HandleUpdateComment)
		})

		// category routes
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", categoryHandler.HandleGetCategories)
		})
	})

	return r
}
