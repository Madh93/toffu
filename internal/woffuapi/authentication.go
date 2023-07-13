package woffuapi

import (
	"net/http"
)

type AuthProvider interface {
	Authenticate(req *http.Request)
	Credentials() string
	Type() string
}

type BasicAuth struct {
	Username string
	Password string
}

func (b *BasicAuth) Authenticate(req *http.Request) {
	req.SetBasicAuth(b.Username, b.Password)
}

func (b *BasicAuth) Credentials() string {
	return b.Username + ":" + b.Password
}

func (b *BasicAuth) Type() string {
	return "BasicAuth"
}

type TokenAuth struct {
	Token string
}

func (t *TokenAuth) Authenticate(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+t.Token)
}

func (t *TokenAuth) Credentials() string {
	return t.Token
}

func (t *TokenAuth) Type() string {
	return "TokenAuth"
}
