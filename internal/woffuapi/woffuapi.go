package woffuapi

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

const (
	BaseURL = "https://app.woffu.com"
)

type APIRequest struct {
	method   string
	endpoint string
	headers  map[string]string
	params   map[string]string
	body     []byte
}

type WoffuAPI struct {
	baseUrl string
	auth    AuthProvider
}

func NewWithToken(token string) *WoffuAPI {
	return &WoffuAPI{
		baseUrl: BaseURL,
		auth:    &TokenAuth{Token: token},
	}
}

func NewWithBasicAuth(username, password string) *WoffuAPI {
	return &WoffuAPI{
		baseUrl: BaseURL,
		auth:    &BasicAuth{Username: username, Password: password},
	}
}

func (w WoffuAPI) makeRequest(request APIRequest) (*http.Response, error) {
	// Build URL
	url, err := w.buildURL(request.endpoint, request.params)
	if err != nil {
		return nil, fmt.Errorf("error building request URL: %v", err)
	}

	// Build basic HTTP request
	req, err := http.NewRequest(request.method, url, bytes.NewBuffer(request.body))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Add HTTP headers
	if request.headers != nil {
		for key, value := range request.headers {
			req.Header.Set(key, value)
		}
	}

	// Setup authentication
	w.auth.Authenticate(req)

	// Do it
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in '%s %s' request: %v", request.method, request.endpoint, err)
	}

	return resp, nil
}

func (w WoffuAPI) buildURL(endpoint string, params map[string]string) (string, error) {
	u, err := url.Parse(w.baseUrl + endpoint)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
