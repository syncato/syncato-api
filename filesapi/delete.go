package filesapi

import (
	"github.com/syncato/syncato-lib/auth"
	"github.com/syncato/syncato-lib/logger"
	"github.com/syncato/syncato-lib/storage/muxstorage"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

func delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	storageMux := ctx.Value("storageMux").(*muxstorage.MuxStorage)
	authRes := ctx.Value("authRes").(*auth.AuthResource)

	rawUri := strings.TrimPrefix(r.URL.Path, "/files/delete/")

	err := storageMux.Remove(authRes, rawUri, true)
	if err != nil {
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
