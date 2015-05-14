package filesapi

import (
	"encoding/json"
	"github.com/syncato/syncato-lib/auth"
	"github.com/syncato/syncato-lib/logger"
	"github.com/syncato/syncato-lib/storage"
	"github.com/syncato/syncato-lib/storage/muxstorage"
	"golang.org/x/net/context"
	"net/http"
	"path/filepath"
)

func move(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*muxstorage.MuxStorage)
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
