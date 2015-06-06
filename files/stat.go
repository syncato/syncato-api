package filesapi

import (
	"encoding/json"
	"github.com/syncato/apis"
	"github.com/syncato/lib/auth"
	"github.com/syncato/lib/logger"
	"github.com/syncato/lib/storage"
	storagemux "github.com/syncato/lib/storage/mux"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"strings"
)

func (api *APIFiles) stat(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*storagemux.StorageMux)
	authRes := ctx.Value("authRes").(*auth.AuthResource)

	rawUri := strings.TrimPrefix(r.URL.RequestURI(), strings.Join([]string{apis.APISROOT, api.GetID(), "stat/"}, "/"))

	var children bool
	queryChildren := r.URL.Query().Get("children")
	if queryChildren != "" {
		ch, err := strconv.ParseBool(queryChildren)
		if err != nil {
			children = false
		}
		children = ch
	}

	meta, err := storageMux.Stat(authRes, rawUri, children)
	if err != nil {
		if storage.IsNotExistError(err) {
			log.Debug(err, nil)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		log.Error("Cannot stat", map[string]interface{}{"err": err})
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(metaJSON)
	return
}
