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
	"io"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	meta "github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/store"
)

type Server struct {
	router *httprouter.Router

	stores   map[string]*store.ServerStore
	storeMtx sync.Mutex
}

func NewServer() *Server {
	s := &Server{
		router: httprouter.New(),
		stores: make(map[string]*store.ServerStore),
	}
	s.router.POST("/:name/tx", s.HandleTransaction)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) HandleTransaction(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var tx meta.SignedTx
	if err := json.NewDecoder(req.Body).Decode(&tx); err != nil {
		http.Error(w, fmt.Sprintf("error decoding request: %s", err), http.StatusBadRequest)
		return
	}
	name := p.ByName("name")
	store, err := s.store(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("error opening store: %s", err), http.StatusInternalServerError)
		return
	}
	hash, err := store.HandleTx(&tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error handling transaction: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, hash.Hex())
}

func (s *Server) store(name string) (*store.ServerStore, error) {
	s.storeMtx.Lock()
	defer s.storeMtx.Unlock()
	if store, ok := s.stores[name]; ok {
		return store, nil
	}
	store, err := store.NewServerStore(name)
	if err != nil {
		return nil, err
	}
	s.stores[name] = store
	return store, nil
}
