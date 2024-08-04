package handlers

import (
	"context"
	"net/http"
	"trip-planner-be/internal/auth"
	"trip-planner-be/internal/config"
	"trip-planner-be/pkg/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func HandleHome(c echo.Context) error {
	return utils.Success(c, map[string]string{"message": "Welcome to the Trip Planner API!"})
}

func HandleLogin(c echo.Context) error {
	url := config.GoogleOAuthConfig.AuthCodeURL(config.OAuthStateString, oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != config.OAuthStateString {
		return utils.Fail(c, map[string]string{"error": "Invalid OAuth state"})
	}

	code := c.QueryParam("code")
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return utils.Error(c, "Failed to exchange token", http.StatusInternalServerError, nil)
	}

	userInfo, err := auth.GetUserInfo(token.AccessToken)
	if err != nil {
		return utils.Error(c, "Failed to get user info", http.StatusInternalServerError, nil)
	}

	combinedInfo := auth.CombineInfo(token, userInfo)
	return utils.Success(c, combinedInfo)
}
