package auth

import (
	"encoding/json"

	"io"
	"net/http"

	"golang.org/x/oauth2"
)

func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func CombineInfo(token *oauth2.Token, userInfo map[string]interface{}) map[string]interface{} {
	tokenInfo := map[string]interface{}{
		"access_token":  token.AccessToken,
		"token_type":    token.TokenType,
		"refresh_token": token.RefreshToken,
		"expiry":        token.Expiry.Unix(),
	}

	return map[string]interface{}{
		"token_info": tokenInfo,
		"user_info":  userInfo,
	}
}
