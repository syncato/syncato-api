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
	"strings"
)

func (api *APIFiles) mkcol(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*storagemux.StorageMux)
	authRes := ctx.Value("authRes").(*auth.AuthResource)
	rawUri := strings.TrimPrefix(r.URL.Path, strings.Join([]string{apis.APISROOT, api.GetID(), "mkcol/"}, "/"))

	err := storageMux.CreateCol(authRes, rawUri, false)
	if err != nil {
		if storage.IsNotExistError(err) {
			log.Debug(err, nil)
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
		if storage.IsExistError(err) {
			log.Debug(err, nil)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	meta, err := storageMux.Stat(authRes, rawUri, false)
	if err != nil {
		if storage.IsNotExistError(err) {
			log.Debug(err, nil)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(metaJSON)
}
