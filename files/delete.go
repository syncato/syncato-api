package filesapi

import (
	"github.com/syncato/apis"
	"github.com/syncato/lib/auth"
	"github.com/syncato/lib/logger"
	storagemux "github.com/syncato/lib/storage/mux"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

func (api *APIFiles) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*storagemux.StorageMux)
	authRes := ctx.Value("authRes").(*auth.AuthResource)

	rawUri := strings.TrimPrefix(r.URL.Path, strings.Join([]string{apis.APISROOT, api.GetID(), "delete/"}, "/"))

	err := storageMux.Remove(authRes, rawUri, true)
	if err != nil {
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
