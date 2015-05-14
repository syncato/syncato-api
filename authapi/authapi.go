package authapi

import (
	"github.com/gorilla/mux"
	"github.com/syncato/syncato-lib/auth/muxauth"
	"github.com/syncato/syncato-lib/config"
	"github.com/syncato/syncato-lib/logger"
	"github.com/syncato/syncato-lib/storage/muxstorage"
)

type API struct {
	router     *mux.Router
	cp         *config.ConfigProvider
	authMux    *muxauth.MuxAuth
	storageMux *muxstorage.MuxStorage
	log        *logger.Logger
}

func NewAPI(router *mux.Router, cp *config.ConfigProvider, authMux *muxauth.MuxAuth, storageMux *muxstorage.MuxStorage, log *logger.Logger) (*API, error) {
	api := API{router, cp, authMux, storageMux, log}
	return &api, nil
}

func (api *API) RegisterRoutes() {
	api.router.HandleFunc("/auth/token/", api.login)
}
