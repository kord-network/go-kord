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

package meta

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	swarmapi "github.com/ethereum/go-ethereum/swarm/api"
	"github.com/meta-network/go-meta/api"
	"github.com/meta-network/go-meta/dapp"
)

type Server struct {
	mux *http.ServeMux

	dapp    *dapp.Dapp
	dappMtx sync.RWMutex

	swarm *swarmapi.Api
}

func NewServer(api *api.API, swarm *swarmapi.Api) *Server {
	s := &Server{
		mux:   http.NewServeMux(),
		swarm: swarm,
	}
	s.mux.Handle("/api/graphql", api)
	s.mux.HandleFunc("/", s.ServeDapp)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
	s.mux.ServeHTTP(w, r)
}

func (s *Server) ServeDapp(w http.ResponseWriter, r *http.Request) {
	s.dappMtx.RLock()
	dapp := s.dapp
	s.dappMtx.RUnlock()

	if dapp == nil {
		http.NotFound(w, r)
		return
	}

	// ensure the root path has a trailing slash so that relative URLs work
	if r.URL.Path == "" {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		return
	}

	key := common.Hex2Bytes(dapp.ManifestHash)
	path := strings.TrimLeft(r.URL.Path, "/")
	reader, contentType, status, err := s.swarm.Get(key, path)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// check the root chunk exists by retrieving the file's size
	if _, err := reader.Size(nil); err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, "", time.Now(), reader)
}

func (s *Server) setDapp(dapp *dapp.Dapp) {
	s.dappMtx.Lock()
	s.dapp = dapp
	s.dappMtx.Unlock()
}
