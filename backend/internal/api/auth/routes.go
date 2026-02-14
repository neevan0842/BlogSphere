package auth

import (
	"fmt"
	"net/http"

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

func (h *handler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	oauthState := h.service.generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(oauthState)
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"url": url,
	})
}

func (h *handler) HandleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, err := r.Cookie("oauthstate")

	if err != nil {
		h.logger.Errorf("could not read oauthstate cookie: %s", err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not authenticate with google"))
		return
	}

	if r.FormValue("state") != oauthState.Value {
		h.logger.Error("invalid oauth google state")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not authenticate with google"))
		return
	}

	// Convert to createUserParams
	userData, err := h.service.getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not authenticate with google"))
		return
	}

	user, isNewUser, err := h.service.createUserIfNotExists(r.Context(), userData)
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

	// Clear the oauthstate cookie
	utils.SetCookie(w, "oauthstate", "", -1)

	// send welcome email
	if isNewUser && user.Username.Valid && user.Email != "" {
		go utils.SendWelcomeEmail(user.Email, user.Username.String, h.logger)
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    fmt.Sprintf("%d minutes", config.Envs.ACCESS_TOKEN_EXPIRE_MINUTES),
	})
}

func (h *handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var payload RefreshRequest

	//parse JSON body
	if err := utils.ParseJSON(r, &payload); err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	// validate refresh token
	if err := utils.Validate.Struct(payload); err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	// get userID from refresh token
	userIDUUID, err := utils.GetUserIDFromToken(w, payload.RefreshToken)
	if err != nil {
		h.logger.Errorf("invalid userID in refresh token: %s", err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
		return
	}

	// check if user exists in database
	_, err = h.service.GetUserByID(r.Context(), userIDUUID)
	if err != nil {
		h.logger.Errorf("failed to get user from database: %s", err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
		return
	}
	// Generate new access token
	accessToken, _, err := utils.GetAccessAndRefreshTokens(userIDUUID.String())
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("could not generate access token"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   fmt.Sprintf("%d minutes", config.Envs.ACCESS_TOKEN_EXPIRE_MINUTES),
	})
}
