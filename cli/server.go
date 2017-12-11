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

package cli

import (
	"net/http"

	swarmhttp "github.com/ethereum/go-ethereum/swarm/api/http"
	"github.com/julienschmidt/httprouter"
	meta "github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/media"
)

type Server struct {
	handler http.Handler
}

func NewServer(store *meta.Store, identityIndex *identity.Index, mediaIndex *media.Index) (*Server, error) {
	srv := &Server{}

	router := httprouter.New()

	identityAPI, err := identity.NewAPI(identityIndex)
	if err != nil {
		return nil, err
	}
	router.Handler("GET", "/meta-id/*path", http.StripPrefix("/meta-id", identityAPI))
	router.Handler("POST", "/meta-id/*path", http.StripPrefix("/meta-id", identityAPI))

	mediaAPI, err := media.NewAPI(mediaIndex, identityIndex)
	if err != nil {
		return nil, err
	}
	router.Handler("GET", "/media/*path", http.StripPrefix("/media", mediaAPI))
	router.Handler("POST", "/media/*path", http.StripPrefix("/media", mediaAPI))

	// add the Swarm API at /bzz: and /bzzr:
	swarmSrv := swarmhttp.NewServer(store.SwarmAPI())
	mux := http.NewServeMux()
	mux.Handle("/bzz:/", swarmSrv)
	mux.Handle("/bzzr:/", swarmSrv)
	mux.Handle("/", router)

	srv.handler = mux

	return srv, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.handler.ServeHTTP(w, req)
}
