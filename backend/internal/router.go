package internal

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	envs "github.com/neevan0842/BlogSphere/backend/config"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/internal/api/auth"

	// mw "github.com/neevan0842/BlogSphere/backend/internal/middleware"
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

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	//
	authService := auth.NewService(sqlc.New(app.db), app.db)
	authHandler := auth.NewHandler(authService, app.logger)

	// authMiddleware := mw.NewMiddleware(sqlc.New(app.db), app.logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/google", authHandler.HandleGoogleLogin)
			r.Get("/google/callback", authHandler.HandleGoogleAuthCallback)
			r.Get("/logout", authHandler.HandleLogout)
			r.Get("/refresh", authHandler.HandleRefresh)
		})
	})

	return r
}
