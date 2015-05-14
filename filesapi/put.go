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

func put(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*muxstorage.MuxStorage)
	authRes := ctx.Value("authRes").(*auth.AuthResource)
	rawUri := strings.TrimPrefix(r.URL.Path, "/files/put/")

	if r.Header.Get("Content-Range") != "" {
		log.Warn("Content-Range not accepted on PUTS", nil)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}

	meta, err := storageMux.Stat(authRes, rawUri, false)
	if err != nil {
		// stat will fail if the file does not exists
		// in our case this is ok and we create a new file
		if storage.IsNotExistError(err) {
			err = storageMux.PutFile(authRes, rawUri, r.Body, r.ContentLength)
			if err != nil {
				log.Error(err, nil)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			meta, err = storageMux.Stat(authRes, rawUri, false)
			if err != nil {
				log.Debug(err, nil)
				if storage.IsNotExistError(err) {
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
			return

		} else {
			log.Error(err, nil)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if meta.IsCol {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	err = storageMux.PutFile(authRes, rawUri, r.Body, r.ContentLength)
	if err != nil {
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	meta, err = storageMux.Stat(authRes, rawUri, false)
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
	return
}
