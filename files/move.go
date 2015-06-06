package filesapi

import (
	"encoding/json"
	"github.com/syncato/lib/auth"
	"github.com/syncato/lib/logger"
	"github.com/syncato/lib/storage"
	storagemux "github.com/syncato/lib/storage/mux"
	"golang.org/x/net/context"
	"net/http"
	"path/filepath"
)

func (api *APIFiles) move(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*storagemux.StorageMux)
	authRes := ctx.Value("authRes").(*auth.AuthResource)

	from := filepath.Clean(r.URL.Query().Get("from"))
	to := filepath.Clean(r.URL.Query().Get("to"))

	err := storageMux.Rename(authRes, from, to)
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

	meta, err := storageMux.Stat(authRes, to, false)
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

	w.WriteHeader(http.StatusOK)
	w.Write(metaJSON)
	return
}
