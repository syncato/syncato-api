package filesapi

import (
	"github.com/syncato/apis"
	authmux "github.com/syncato/lib/auth/mux"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type APIFiles struct {
	id string
}

func NewAPIFiles(id string) *APIFiles {
	fapi := APIFiles{
		id: id,
	}
	return &fapi
}
func (api *APIFiles) GetID() string { return api.id }

func (api *APIFiles) HandleRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	authMux := ctx.Value("authMux").(*authmux.AuthMux)
	authMid := authMux.AuthMiddleware
	path := r.URL.Path

	if strings.HasPrefix(path, strings.Join([]string{apis.APISROOT, api.GetID(), "get"}, "/")) && r.Method == "GET" {
		authMid(ctx, w, r, api.get)
	} else if strings.HasPrefix(path, strings.Join([]string{apis.APISROOT, api.GetID(), "delete"}, "/")) && r.Method == "DELETE" {
		authMid(ctx, w, r, api.delete)
	} else if strings.HasPrefix(path, strings.Join([]string{apis.APISROOT, api.GetID(), "mkcol"}, "/")) && r.Method == "POST" {
		authMid(ctx, w, r, api.mkcol)
	} else if strings.HasPrefix(path, strings.Join([]string{apis.APISROOT, api.GetID(), "move"}, "/")) && r.Method == "POST" {
		authMid(ctx, w, r, api.move)
	} else if strings.HasPrefix(path, strings.Join([]string{apis.APISROOT, api.GetID(), "put"}, "/")) && r.Method == "PUT" {
		authMid(ctx, w, r, api.put)
	} else if strings.HasPrefix(path, strings.Join([]string{apis.APISROOT, api.GetID(), "stat"}, "/")) && r.Method == "GET" {
		authMid(ctx, w, r, api.stat)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}
