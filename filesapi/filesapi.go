package filesapi

import (
	"fmt"
	"github.com/syncato/syncato-lib/auth/muxauth"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

func myAction(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, http.ResponseWriter, *http.Request, bool) {
	fmt.Println("hyujuuuu")
	return ctx, w, r, false
}

type FilesAPIOp struct {
	id     string
	method string
	action func(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

func (op *FilesAPIOp) GetID() string     { return op.id }
func (op *FilesAPIOp) GetMethod() string { return op.method }
func (op *FilesAPIOp) Execute(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_ctx, _w, _r, shouldContinue := myAction(ctx, w, r)
	if !shouldContinue {
		op.action(_ctx, _w, _r)
	}
}

type FilesAPI struct {
	id  string
	ops []*FilesAPIOp
}

func (api *FilesAPI) GetID() string         { return api.id }
func (api *FilesAPI) GetOps() []*FilesAPIOp { return api.ops }

func NewFilesAPI() *FilesAPI {
	fapi := FilesAPI{
		id: "files",
		ops: []*FilesAPIOp{
			&FilesAPIOp{"get", "GET", get},
			&FilesAPIOp{"delete", "DELETE", delete},
			&FilesAPIOp{"mkcol", "POST", mkcol},
			&FilesAPIOp{"move", "POST", move},
			&FilesAPIOp{"put", "PUT", put},
			&FilesAPIOp{"stat", "STAT", stat},
		},
	}
	return &fapi
}

func HandleRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	authMux := ctx.Value("authMux").(*muxauth.MuxAuth)
	authMid := authMux.AuthMiddleware
	path := r.URL.Path

	if strings.HasPrefix(path, "/api/files/get/") && r.Method == "GET" {
		authMid(ctx, w, r, get)
	} else if strings.HasPrefix(path, "/api/files/delete/") && r.Method == "DELETE" {
		authMid(ctx, w, r, delete)
	} else if strings.HasPrefix(path, "/api/files/mkcol/") && r.Method == "POST" {
		authMid(ctx, w, r, mkcol)
	} else if strings.HasPrefix(path, "/api/files/move/") && r.Method == "POST" {
		authMid(ctx, w, r, move)
	} else if strings.HasPrefix(path, "/api/files/put/") && r.Method == "PUT" {
		authMid(ctx, w, r, put)
	} else if strings.HasPrefix(path, "/api/files/stat/") && r.Method == "GET" {
		authMid(ctx, w, r, get)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}
