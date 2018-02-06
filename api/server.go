// This file is part of the go-meta library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/cayleygraph/cayley/graph"
	"github.com/julienschmidt/httprouter"
	"github.com/meta-network/go-meta/db"
)

// Server implements the META graph API which supports creating and updating
// META graphs.
type Server struct {
	router *httprouter.Router

	stores   map[string]graph.QuadStore
	storeMtx sync.Mutex
}

// NewServer returns a new server.
func NewServer() *Server {
	s := &Server{
		router: httprouter.New(),
		stores: make(map[string]graph.QuadStore),
	}
	s.router.POST("/:name", s.HandleCreate)
	s.router.POST("/:name/apply-deltas", s.HandleApplyDeltas)
	return s
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

// HandleCreate handles a request to create a META graph.
func (s *Server) HandleCreate(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	if err := db.Create(name); err != nil {
		http.Error(w, fmt.Sprintf("error creating database: %s", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleApplyDeltas handles a request to update a META graph.
func (s *Server) HandleApplyDeltas(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var r Request
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		http.Error(w, fmt.Sprintf("error decoding request: %s", err), http.StatusBadRequest)
		return
	}
	name := p.ByName("name")
	store, err := s.store(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("error opening store: %s", err), http.StatusInternalServerError)
		return
	}
	if err := store.ApplyDeltas(r.In, r.Opts); err != nil {
		http.Error(w, fmt.Sprintf("error applying deltas: %s", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) store(name string) (graph.QuadStore, error) {
	s.storeMtx.Lock()
	defer s.storeMtx.Unlock()
	if store, ok := s.stores[name]; ok {
		return store, nil
	}
	store, err := graph.NewQuadStore("meta", name, graph.Options{})
	if err != nil {
		return nil, err
	}
	s.stores[name] = store
	return store, nil
}
