package filesapi

import (
	"encoding/json"
	"github.com/syncato/syncato-lib/auth"
	"github.com/syncato/syncato-lib/logger"
	"github.com/syncato/syncato-lib/storage"
	"github.com/syncato/syncato-lib/storage/muxstorage"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

func mkcol(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*muxstorage.MuxStorage)
	authRes := ctx.Value("authRes").(*auth.AuthResource)
	rawUri := strings.TrimPrefix(r.URL.Path, "/files/mkcol/")

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
