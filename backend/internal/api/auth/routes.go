package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/neevan0842/BlogSphere/backend/config"
	"github.com/neevan0842/BlogSphere/backend/utils"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type handler struct {
	service Service
	logger  *zap.SugaredLogger
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     config.Envs.GOOGLE_CLIENT_ID,
	ClientSecret: config.Envs.GOOGLE_CLIENT_SECRET,
	RedirectURL:  config.Envs.GOOGLE_REDIRECT_URI,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func NewHandler(service Service, logger *zap.SugaredLogger) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/google", h.handleGoogleLogin)
		r.Get("/google/callback", h.handleGoogleAuthCallback)
		// TODO: add auth middleware to protect this route
		r.Get("/refresh", h.handleRefresh)
		r.Get("/logout", h.handleLogout)
	})
}

func (h *handler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	oauthState := h.service.generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) handleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, err := r.Cookie("oauthstate")

	if err != nil {
		h.logger.Error("could not read oauthstate cookie: %s", err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not authenticate with google"))
		return
	}

	if r.FormValue("state") != oauthState.Value {
		h.logger.Error("invalid oauth google state")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not authenticate with google"))
		return
	}

	// Convert to createOrUpdateUserParams
	userData, err := h.service.getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not authenticate with google"))
		return
	}

	user, err := h.service.createOrUpdateUser(context.Background(), userData)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not authenticate with google"))
		return
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := utils.GetAccessAndRefreshTokens(user.ID.String())
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not generate tokens"))
		return
	}
	// Set tokens in cookies
	utils.SetCookie(w, "access_token", accessToken, config.Envs.ACCESS_TOKEN_EXPIRE_MINUTES)
	utils.SetCookie(w, "refresh_token", refreshToken, config.Envs.REFRESH_TOKEN_EXPIRE_MINUTES)
	utils.SetCookie(w, "oauthstate", "", -1) // Clear the oauthstate cookie

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Successfully authenticated with Google"})
}

func (h *handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear the access_token and refresh_token cookies
	utils.SetCookie(w, "access_token", "", -1)
	utils.SetCookie(w, "refresh_token", "", -1)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

func (h *handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	// TODO: implement token refresh logic
}
