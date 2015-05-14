package authapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

func (api *API) login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	params := loginParams{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		api.log.Debug(err, nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	authRes, err := api.authMux.Authenticate(params.Username, params.Password, params.ID)
	if err != nil {
		api.log.Debug(err, nil)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	tokenString, err := api.authMux.CreateAuthTokenFromAuthResource(authRes)
	if err != nil {
		api.log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := make(map[string]string)
	data["auth_token"] = tokenString
	tokenJSON, err := json.Marshal(data)
	if err != nil {
		api.log.Error(err, nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(tokenJSON)
	return
}
