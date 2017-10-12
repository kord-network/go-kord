// This file is part of the go-meta library.
//
// Copyright (C) 2017 JAAK MUSIC LTD
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// If you have any questions please contact yo@jaak.io

package identity

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// API is a HTTP API to create and update META identities.
type API struct {
	store  Store
	router *httprouter.Router
}

// NewAPI returns a new API which uses the given store to load and save
// identities.
func NewAPI(store Store) *API {
	api := &API{
		store:  store,
		router: httprouter.New(),
	}
	api.router.POST("/", api.handleCreateIdentity)
	api.router.GET("/:name", api.handleGetIdentity)
	api.router.POST("/:name", api.handleUpdateIdentity)
	return api
}

// ServeHTTP implements the http.Handler interface so that the API can be used
// to serve HTTP requests.
func (a *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a *API) handleCreateIdentity(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	name := req.FormValue("name")
	if name == "" {
		http.Error(w, "missing name parameter", http.StatusBadRequest)
		return
	}
	identity := NewIdentity(name)
	if err := a.store.Save(identity); err != nil {
		http.Error(w, fmt.Sprintf("error saving identity: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application./json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(identity)
}

func (a *API) handleGetIdentity(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	identity, err := a.store.Load(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("error loading identity: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application./json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(identity)
}

func (a *API) handleUpdateIdentity(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var aux map[string]string
	if err := json.NewDecoder(req.Body).Decode(&aux); err != nil {
		http.Error(w, fmt.Sprintf("error reading aux from request: %s", err), http.StatusBadRequest)
		return
	}
	name := p.ByName("name")
	identity, err := a.store.Load(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("error loading identity: %s", err), http.StatusInternalServerError)
		return
	}
	for key, val := range aux {
		identity.Aux[key] = val
	}
	if err := a.store.Save(identity); err != nil {
		http.Error(w, fmt.Sprintf("error saving identity: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application./json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(identity)
}
