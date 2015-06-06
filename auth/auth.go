package auth

import (
	"github.com/syncato/apis"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type APIAuth struct {
	id string
}

func NewAPIAuth(id string) *APIAuth {
	fapi := APIAuth{
		id: id,
	}
	return &fapi
}
func (api *APIAuth) GetID() string { return api.id }

func (api *APIAuth) HandleRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasPrefix(path, strings.Join([]string{apis.APISROOT, api.GetID(), "login"}, "/")) && r.Method == "POST" {
		api.login(ctx, w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}
