package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"trip-planner-be/internal/config"

	"golang.org/x/oauth2"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/login">Google Log In</a></body></html>`
	fmt.Fprint(w, html)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOAuthConfig.AuthCodeURL(config.OAuthStateString,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if token.AccessToken == "" {
		fmt.Println("Error: Access token is empty")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Print token information
	fmt.Printf("Access Token: %s\n", token.AccessToken)
	fmt.Printf("Token Type: %s\n", token.TokenType)
	fmt.Printf("Refresh Token: %s\n", token.RefreshToken)
	fmt.Printf("Expiry: %v\n", token.Expiry)

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("Failed getting user info: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed reading response body: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		fmt.Printf("Failed decoding user info: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Combine token info and user info
	tokenInfo := map[string]interface{}{
		"access_token":  token.AccessToken,
		"token_type":    token.TokenType,
		"refresh_token": token.RefreshToken,
		"expiry":        token.Expiry,
	}

	combinedInfo := map[string]interface{}{
		"token_info": tokenInfo,
		"user_info":  userInfo,
	}

	// Convert combined info to JSON
	jsonResponse, err := json.MarshalIndent(combinedInfo, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(jsonResponse)
}
