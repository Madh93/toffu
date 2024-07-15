package woffuapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Token string

type TokenResponse struct {
	AccessToken Token `json:"access_token"`
}

func (w WoffuAPI) CreateToken() (*TokenResponse, error) {
	if w.auth.Type() != "BasicAuth" {
		return nil, errors.New("basic authentication (user/pass) is required")
	}

	credentials := strings.Split(w.auth.Credentials(), ":")

	// Build API Request
	body := fmt.Sprintf("grant_type=password&username=%s&password=%s", credentials[0], credentials[1])
	apiRequest := APIRequest{
		method:   "POST",
		endpoint: "/token",
		headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/x-www-form-urlencoded",
		},
		// Credentials need to be passed in the body :S
		body: []byte(body),
	}

	// Get Token
	resp, err := w.makeRequest(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error creating Token")
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}
