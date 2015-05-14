package filesapi

import (
	"github.com/syncato/syncato-lib/auth"
	"github.com/syncato/syncato-lib/logger"
	"github.com/syncato/syncato-lib/storage"
	"github.com/syncato/syncato-lib/storage/muxstorage"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func get(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*muxstorage.MuxStorage)
	authRes := ctx.Value("authRes").(*auth.AuthResource)

	log.Info("GGGGEEEETTTTINGGG...", nil)
	rawUri := strings.TrimPrefix(r.URL.Path, "/api/files/get/")

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

	if meta.IsCol {
		// TODO: here we could do the zip based download for folders
		log.Warn("GET not implemented for non-cols", nil)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}

	reader, err := storageMux.GetFile(authRes, rawUri)
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
	w.Header().Set("Content-Type", meta.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(meta.Path))
	w.WriteHeader(http.StatusOK)

	io.Copy(w, reader)

	return
}
