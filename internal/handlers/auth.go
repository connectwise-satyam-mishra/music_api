// Package handlers provides HTTP handlers for authentication and authorization.

package handlers

import (
	"context"
	"fmt"
	"music_api/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// spotifyOauthConfig is the OAuth2 configuration for Spotify authentication.
var spotifyOauthConfig *oauth2.Config

// oauthStateString is a random string used for OAuth2 state validation.
var oauthStateString = "random"

// Init sets up the Spotify OAuth2 configuration.
func Init() {
	spotifyOauthConfig = &oauth2.Config{
		ClientID:     config.Config.ClientId,
		ClientSecret: config.Config.ClientSecret,
		RedirectURL:  config.Config.AppURL + "/auth/spotify/callback",
		Scopes:       []string{"user-read-private"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Config.SpotifyAuthURL + "/authorize",
			TokenURL: config.Config.SpotifyAuthURL + "/api/token",
		},
	}
}

// SpotifyAuthMiddleware handles authentication using Spotify OAuth2.
func SpotifyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is authenticated, otherwise initiate Spotify authentication flow
		// Allowing to validate token from cookie or header
		authToken, err := c.Request.Cookie("access_token")
		token := ""
		bearer_token := c.Request.Header.Get("Authorization")
		if bearer_token != "" {
			token = strings.Split(bearer_token, " ")[1]
		} else if err == nil {
			token = authToken.Value
		}

		isValid, err := validateSpotifyToken(token)
		if !isValid || err != nil {
			// Redirect to Spotify authentication URL
			url := spotifyOauthConfig.AuthCodeURL(oauthStateString)
			c.Redirect(http.StatusTemporaryRedirect, url)
			c.Abort()
			return
		}
		c.Next()
	}
}

// validateSpotifyToken validates the Spotify access token.
func validateSpotifyToken(accessToken string) (bool, error) {
	if accessToken == "" {
		return false, nil
	}
	url := config.Config.SpotifyAPIURL + "/v1/me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil // Token is valid
	}

	return false, nil // Token is not valid
}

// CallbackHandler handles the callback after the user authorizes via Spotify.
func CallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		if state != oauthStateString {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid oauth state"))
			return
		}
		code := c.Query("code")
		token, err := spotifyOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Store the access token in a cookie for subsequent requests
		c.SetCookie("access_token", token.AccessToken, 3600, "/", "localhost", false, true)

		// Redirect or perform action after successful authentication
		c.Redirect(http.StatusTemporaryRedirect, "/api/")
	}
}
