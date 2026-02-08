package auth

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/neevan0842/BlogSphere/backend/config"
	"github.com/neevan0842/BlogSphere/backend/logger"
	"github.com/neevan0842/BlogSphere/backend/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type handler struct {
	service Service
}

var log = logger.Get()

var googleOauthConfig = &oauth2.Config{
	ClientID:     config.Envs.GOOGLE_CLIENT_ID,
	ClientSecret: config.Envs.GOOGLE_CLIENT_SECRET,
	RedirectURL:  config.Envs.GOOGLE_REDIRECT_URI,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/google", h.handleGoogleLogin)
		r.Get("/google/callback", h.handleGoogleAuthCallback)
		r.Get("/logout", h.handleLogout)
	})
}

func (h *handler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) handleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Error("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Error(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// TODO: Create or update user in database and generate JWT token and set it in cookie
	// TODO: Redirect to frontend with token in cookie
	fmt.Println(string(data))
	utils.WriteJSON(w, http.StatusOK, data)
}

func (h *handler) handleLogout(w http.ResponseWriter, r *http.Request) {}
