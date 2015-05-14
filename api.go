package api

import (
	"golang.org/x/net/context"
	"net/http"
)

type API interface {
	GetID() string
	GetOps() []APIOp
}

type APIOp interface {
	GetID() string
	GetMethod() string
	Execute(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

type APIAction interface {
	Execute(ctx context.Context, w http.ResponseWriter, r *http.Request) bool
}
