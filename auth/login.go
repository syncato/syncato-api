package auth

import (
	"encoding/json"
	authmux "github.com/syncato/lib/auth/mux"
	"github.com/syncato/lib/config"
	"github.com/syncato/lib/logger"
	storagemux "github.com/syncato/lib/storage/mux"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
	Extra    string `json:extra`
}

func (api *APIAuth) login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log := ctx.Value("log").(*logger.Logger)
	authMux := ctx.Value("authMux").(*authmux.AuthMux)
	storageMux := ctx.Value("storageMux").(*storagemux.StorageMux)
	cfg := ctx.Value("cfg").(*config.Config)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	params := loginParams{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		log.Debug(err, nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	authRes, err := authMux.Authenticate(params.Username, params.Password, params.ID, params.Extra)
	if err != nil {
		log.Debug(err, nil)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// We create the homedir for the user if the config option "createUserHomeOnLogin" is true.
	if cfg.CreateUserHomeOnLogin() == true {
		for _, scheme := range cfg.CreateUserHomeInStorages() {
			ok, err := storageMux.IsUserHomeCreated(authRes, scheme)
			if err != nil {
				log.Error("Checking existence of user home failed", map[string]interface{}{
					"err":            err,
					"auth_id":        authRes.AuthID,
					"username":       authRes.Username,
					"storage_scheme": scheme,
				})
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if !ok {
				// we create the user home
				err := storageMux.CreateUserHome(authRes, scheme)
				if err != nil {
					log.Error("Creation of user home failed", map[string]interface{}{
						"err":            err,
						"auth_id":        authRes.AuthID,
						"username":       authRes.Username,
						"storage_scheme": scheme,
					})
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	tokenString, err := authMux.CreateAuthTokenFromAuthResource(authRes)
	if err != nil {
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := make(map[string]string)
	data["auth_token"] = tokenString
	tokenJSON, err := json.Marshal(data)
	if err != nil {
		log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(tokenJSON)
	return
}
