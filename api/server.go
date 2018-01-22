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

	"github.com/julienschmidt/httprouter"
	meta "github.com/meta-network/go-meta"
)

type Server struct {
	state  *meta.State
	router *httprouter.Router
}

func NewServer(state *meta.State) *Server {
	s := &Server{
		state:  state,
		router: httprouter.New(),
	}
	s.router.POST("/tx", s.HandleTransaction)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) HandleTransaction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var tx meta.Tx
	if err := json.NewDecoder(req.Body).Decode(&tx); err != nil {
		http.Error(w, fmt.Sprintf("error decoding request: %s", err), http.StatusBadRequest)
		return
	}
	hash, err := s.state.Apply(&tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error applying transaction: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, hash.Hex())
}
