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

package stream

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flynn/flynn/pkg/sse"
	"github.com/julienschmidt/httprouter"
	"github.com/meta-network/go-meta"
	"gopkg.in/inconshreveable/log15.v2"
)

// API is a HTTP API to read from and write to META streams.
type API struct {
	store  *meta.Store
	router *httprouter.Router
}

// NewAPI returns a new API which gets streams from the given store.
func NewAPI(store *meta.Store) *API {
	api := &API{
		store:  store,
		router: httprouter.New(),
	}
	api.router.GET("/:name", api.handleReadStream)
	api.router.POST("/:name", api.handleWriteStream)
	return api
}

// ServeHTTP implements the http.Handler interface so that the API can be used
// to serve HTTP requests.
func (a *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a *API) handleReadStream(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	reader, err := a.store.StreamReader(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting stream reader: %s", err), http.StatusInternalServerError)
		return
	}
	defer reader.Close()
	ch := make(chan *meta.Object)
	stream := sse.NewStream(w, ch, log15.New())
	go func() {
		for {
			id, ok := <-reader.Ch()
			if !ok {
				if err := reader.Err(); err != nil {
					stream.CloseWithError(err)
				}
				return
			}
			obj, err := a.store.Get(id)
			if err != nil {
				stream.CloseWithError(err)
				return
			}
			select {
			case ch <- obj:
			case <-stream.Done:
				return
			}
		}
	}()
	stream.Serve()
	stream.Wait()
}

func (a *API) handleWriteStream(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	writer, err := a.store.StreamWriter(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting stream writer: %s", err), http.StatusInternalServerError)
		return
	}
	defer writer.Close()
	var v map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&v); err != nil {
		http.Error(w, fmt.Sprintf("error reading body: %s", err), http.StatusBadRequest)
		return
	}
	obj, err := a.store.Put(v)
	if err != nil {
		http.Error(w, fmt.Sprintf("error storing object: %s", err), http.StatusInternalServerError)
		return
	}
	if err := writer.Write(obj.Cid()); err != nil {
		http.Error(w, fmt.Sprintf("error writing object to stream: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj.Cid())
}
