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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	swarmhttp "github.com/ethereum/go-ethereum/swarm/api/http"
	"github.com/ipfs/go-cid"
	"github.com/julienschmidt/httprouter"
	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/cwr"
	"github.com/meta-network/go-meta/eidr"
	"github.com/meta-network/go-meta/ern"
	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/media"
	"github.com/meta-network/go-meta/musicbrainz"
	"github.com/meta-network/go-meta/stream"
	"github.com/meta-network/go-meta/xml"
)

type Server struct {
	handler http.Handler
	store   *meta.Store
}

func NewServer(store *meta.Store, indexes map[string]*meta.Index) (*Server, error) {
	srv := &Server{store: store}

	router := httprouter.New()
	router.GET("/object/:cid", srv.HandleGetObject)
	router.POST("/convert/xml", srv.HandleConvertXML)

	// add the identity API at /meta-id
	identityStore := identity.NewMemoryStore()
	identityAPI := identity.NewAPI(identityStore)
	router.Handler("GET", "/meta-id/*path", http.StripPrefix("/meta-id", identityAPI))
	router.Handler("POST", "/meta-id/*path", http.StripPrefix("/meta-id", identityAPI))

	// add the stream API at /stream
	streamAPI := stream.NewAPI(store)
	router.Handler("GET", "/stream/*path", http.StripPrefix("/stream", streamAPI))
	router.Handler("POST", "/stream/*path", http.StripPrefix("/stream", streamAPI))

	mediaResolver := &media.Resolver{
		Store:   store,
		IDStore: identityStore,
	}

	if index, ok := indexes["musicbrainz"]; ok {
		api, err := musicbrainz.NewAPI(index.DB, store)
		if err != nil {
			return nil, err
		}
		mediaResolver.MusicBrainz = api.Resolver()
		router.Handler("GET", "/musicbrainz/*path", http.StripPrefix("/musicbrainz", api))
		router.Handler("POST", "/musicbrainz/*path", http.StripPrefix("/musicbrainz", api))
	}

	if index, ok := indexes["cwr"]; ok {
		api, err := cwr.NewAPI(index.DB, store)
		if err != nil {
			return nil, err
		}
		mediaResolver.Cwr = api.Resolver()
		router.Handler("GET", "/cwr/*path", http.StripPrefix("/cwr", api))
		router.Handler("POST", "/cwr/*path", http.StripPrefix("/cwr", api))
	}

	if index, ok := indexes["ern"]; ok {
		api, err := ern.NewAPI(index.DB, store)
		if err != nil {
			return nil, err
		}
		mediaResolver.Ern = api.Resolver()
		router.Handler("GET", "/ern/*path", http.StripPrefix("/ern", api))
		router.Handler("POST", "/ern/*path", http.StripPrefix("/ern", api))
	}
	if index, ok := indexes["eidr"]; ok {
		api, err := eidr.NewAPI(index.DB, store)
		if err != nil {
			return nil, err
		}
		router.Handler("GET", "/eidr/*path", http.StripPrefix("/eidr", api))
		router.Handler("POST", "/eidr/*path", http.StripPrefix("/eidr", api))
	}

	mediaAPI, err := media.NewAPI(mediaResolver)
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

func (s *Server) HandleConvertXML(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	source := req.URL.Query().Get("source")
	if source == "" {
		http.Error(w, "missing source parameter", http.StatusBadRequest)
		return
	}

	var context []*cid.Cid
	if c := req.URL.Query().Get("context"); c != "" {
		for _, v := range strings.Split(c, ",") {
			cid, err := cid.Decode(v)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid CID in context value %q: %s", v, err), http.StatusBadRequest)
				return
			}
			context = append(context, cid)
		}
	}

	converter := metaxml.NewConverter(s.store)
	obj, err := converter.ConvertXML(req.Body, context, source)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj)
}

func (s *Server) HandleGetObject(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	cid, err := cid.Decode(p.ByName("cid"))
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid CID %q: %s", p.ByName("cid"), err), http.StatusBadRequest)
		return
	}

	obj, err := s.store.Get(cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj)
}
